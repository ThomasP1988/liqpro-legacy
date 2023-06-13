package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"

	coinbase "liqpro/exchanges/coinbase"
	disruptor "liqpro/shared/disruptor/parser"
)

var mapChannels map[string]string = map[string]string{
	"BTC-EUR": "btceur",
	"ETH-EUR": "etheur",
}
var markets = []string{"BTC-EUR", "ETH-EUR"}
var orderbookSnapshot = map[string]bool{}
var orderbookSnapshotLastUpdateID = &map[string]int64{}

var lastProcessedEventMicroTimestamp int64 = 0

var asks = map[string]map[string]string{}
var bids = map[string]map[string]string{}

var mutexMarkets *disruptor.MutexMarkets

var startUpdateByte = []byte("{\"type\":\"l2update")
var startSnapshotByte = []byte("{\"type\":\"snapshot")

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	fees, err := coinbase.GetFees()

	if err != nil {
		println("GetFees", err)
		panic(err)
	}

	channels := []string{}

	for _, value := range mapChannels {
		channels = append(channels, value)
		asks[value] = map[string]string{}
		bids[value] = map[string]string{}
		orderbookSnapshot[value] = false
	}

	mutexMarkets = disruptor.SetMutexes(&channels)

	wsDepthHandler := disruptor.Disruptor(NewConsumer(fees))

	go coinbase.WsDepthServe(markets, coinbase.WsDepthHandler(wsDepthHandler))

	for {
		select {
		case <-interrupt:
			fmt.Println("interrupt")
			return
		}
	}
}

func parser(msg *[]byte, csm *Consumer) {
	println(string(*msg))
	if bytes.HasPrefix(*msg, startUpdateByte) {
		csm.err = json.Unmarshal(*msg, csm.event)
		if csm.err != nil {
			fmt.Println("error unmarshalling", string(*msg), csm.err)
		} else {
			ProcessEvent(csm)
		}
		fmt.Println(orderbookSnapshot)
	} else if bytes.HasPrefix(*msg, startSnapshotByte) {
		csm.err = json.Unmarshal(*msg, csm.snap)
		if csm.err != nil {
			fmt.Println("error unmarshalling", string(*msg), csm.err)
		}
		fmt.Println("SNAPSHOT", csm.snap.ProductID)
		ProcessSnap(csm)

		csm.boolHelper = true
		for _, value := range orderbookSnapshot {
			if !value {
				csm.boolHelper = false
			}

		}

		if csm.boolHelper {
			debug.SetGCPercent(-1)
		}
	} else if bytes.HasPrefix(*msg, []byte("{\"type\":\"error")) {
		fmt.Println("error", string(*msg))
	}

}

// ProcessEvent process event to populate map and DB
func ProcessEvent(csm *Consumer) {
	csm.B.Market = mapChannels[csm.event.ProductID]
	(*mutexMarkets)[csm.B.Market].Lock()
	defer (*mutexMarkets)[csm.B.Market].Unlock()
	csm.dataLen = len(csm.event.Changes)

	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {
		if csm.event.Changes[csm.i][0] == "sell" {
			csm.B.Action = "1"
			if csm.previousVolume, csm.ok = asks[csm.B.Market][csm.event.Changes[csm.i][1]]; csm.ok {

				Del(csm, &csm.event.Changes[csm.i][1], &csm.previousVolume)
			}
			if csm.event.Changes[csm.i][2] != "0.00000000" {
				Upsert(csm, &csm.event.Changes[csm.i][1], &csm.event.Changes[csm.i][2])
				asks[csm.B.Market][csm.event.Changes[csm.i][1]] = csm.event.Changes[csm.i][2]
			} else {
				delete(asks[csm.B.Market], csm.event.Changes[csm.i][1])
			}
		} else {
			csm.B.Action = "0"
			if csm.previousVolume, csm.ok = bids[csm.B.Market][csm.event.Changes[csm.i][1]]; csm.ok {
				Del(csm, &csm.event.Changes[csm.i][1], &csm.previousVolume)
			}
			if csm.event.Changes[csm.i][2] != "0.00000000" {
				Upsert(csm, &csm.event.Changes[csm.i][1], &csm.event.Changes[csm.i][2])
				bids[csm.B.Market][csm.event.Changes[csm.i][1]] = csm.event.Changes[csm.i][2]
			} else {
				delete(bids[csm.B.Market], csm.event.Changes[csm.i][1])
			}
		}
	}

}

// ProcessSnap first event is a snapshot
func ProcessSnap(csm *Consumer) {
	csm.B.Market = mapChannels[csm.snap.ProductID]
	(*mutexMarkets)[csm.B.Market].Lock()
	defer (*mutexMarkets)[csm.B.Market].Unlock()

	csm.dataLen = len(csm.snap.Asks)
	// if csm.dataLen > 50 {
	// 	csm.dataLen = 50
	// }
	csm.B.Action = "1"
	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {
		Upsert(csm, &csm.snap.Asks[csm.i][0], &csm.snap.Asks[csm.i][1])
		asks[csm.B.Market][csm.snap.Asks[csm.i][0]] = csm.snap.Asks[csm.i][1]
	}
	csm.dataLen = len(csm.snap.Bids)
	// if csm.dataLen > 50 {
	// 	csm.dataLen = 50
	// }
	csm.B.Action = "0"
	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {
		Upsert(csm, &csm.snap.Bids[csm.i][0], &csm.snap.Bids[csm.i][1])
		bids[csm.B.Market][csm.snap.Bids[csm.i][0]] = csm.snap.Bids[csm.i][1]
	}
	orderbookSnapshot[csm.B.Market] = true

}

func Del(csm *Consumer, price *string, volume *string) {
	// 000000

	csm.B.Price = (*price)
	csm.B.Volume = strings.Replace(*volume, ".", "", 1)

	disruptor.Del(&csm.B)

}
func Upsert(csm *Consumer, price *string, volume *string) {
	// 000000
	csm.B.Price = (*price)
	csm.B.Volume = strings.Replace(*volume, ".", "", 1)
	generateScore(csm)
	disruptor.Upsert(&csm.B)
}

func generateScore(csm *Consumer) error {
	csm.tmpFloat, csm.err = strconv.ParseFloat(csm.B.Price, 64)
	if csm.err != nil {
		return csm.err
	}
	csm.B.Score = strings.Replace(strconv.FormatFloat(csm.tmpFloat+(csm.tmpFloat*csm.feeTaker), 'f', 2, 64), ".", "", 1)
	return nil
}
