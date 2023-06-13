package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"liqpro/exchanges/bitstamp"
	disruptor "liqpro/shared/disruptor/parser"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
)

var prefixSubscriptionSucceeded []byte = []byte("{\"event\":\"bts:subscription_succeeded")

var balance *bitstamp.AccountBalanceResult
var err error

var mapChannels map[string]string = map[string]string{
	"order_book_etheur": "etheur",
	"order_book_btceur": "btceur",
	"order_book_btcusd": "btcusd",
}

var markets = []string{"btceur", "btcusd", "etheur"}

var mapFees map[string]float64

var orderbookSnapshot = make(map[string]bool, len(markets))
var orderbookSnapshotMicrotimestamp = make(map[string]int64, len(markets))

var microtimestampByMarket = make(map[string]int64, len(markets))

var lastProcessedEventMicroTimestamp int64 = 0

var asks = make(map[string]*map[string]string, len(markets))
var bids = make(map[string]*map[string]string, len(markets))

var mutexMarkets *disruptor.MutexMarkets
var eventQueue *EventQueue

var FirstConsumer = NewConsumer()

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	mutexMarkets = disruptor.SetMutexes(&markets)
	eventQueue = SetEventQueue(&markets)

	balance, err = bitstamp.GetAccount()
	if err != nil {
		panic(err)
	}
	println("balance", *&balance.BtcEurFee)

	mapFees = map[string]float64{
		"etheur": balance.EthEurFee / 100,
		"btceur": balance.BtcEurFee / 100,
		"btcusd": balance.BtcUsdFee / 100,
	}

	wsDepthHandler := disruptor.Disruptor(
		FirstConsumer,
		// NewConsumer(),
	)

	var onConnected = func() {
		debug.SetGCPercent(-1)
	}

	go bitstamp.WsDepthServe(markets, bitstamp.WsDepthHandler(wsDepthHandler), &onConnected)

	for _, market := range markets {
		orderbookSnapshot[market] = false
		asksMap := make(map[string]string, 30)
		bidsMap := make(map[string]string, 30)
		asks[market] = &asksMap
		bids[market] = &bidsMap
	}

	for {
		select {
		case <-interrupt:
			fmt.Println("interrupt")
			return
		}
	}

}

func parser(msg *[]byte, csm *Consumer) {
	println("msg", string(*msg))

	if bytes.HasPrefix(*msg, prefixSubscriptionSucceeded) {
		return
	}

	csm.err = json.Unmarshal(*msg, csm.event)
	if csm.err != nil {
		fmt.Println("error unmarshalling", string(*msg), csm.err)
	}
	csm.B.Market = mapChannels[csm.event.Channel]

	ProcessEvent(csm)

}

// ProcessEvent process event to populate map and DB
func ProcessEvent(csm *Consumer) {
	(*mutexMarkets)[csm.B.Market].Lock()
	csm.B.Action = "1"
	for csm.i = 0; csm.i < 25; csm.i++ {
		if csm.previousVolume, csm.ok = (*asks[csm.B.Market])[csm.event.Data.Asks[csm.i][0]]; csm.ok {

			if csm.event.Data.Asks[csm.i][1] == (*asks[csm.B.Market])[csm.event.Data.Asks[csm.i][0]] {
				delete((*asks[csm.B.Market]), csm.event.Data.Asks[csm.i][0])
				(*csm.newAsks)[csm.event.Data.Asks[csm.i][0]] = csm.event.Data.Asks[csm.i][1]
				continue
			}

			delete((*asks[csm.B.Market]), csm.event.Data.Asks[csm.i][0])
			csm.B.Price = csm.event.Data.Asks[csm.i][0]
			Del(csm, &csm.previousVolume)
		}

		csm.B.Price = csm.event.Data.Asks[csm.i][0]

		Upsert(csm, &csm.event.Data.Asks[csm.i][1])
		(*csm.newAsks)[csm.event.Data.Asks[csm.i][0]] = csm.event.Data.Asks[csm.i][1]
	}

	csm.B.Action = "0"
	for csm.i = 0; csm.i < 25; csm.i++ {
		if csm.previousVolume, csm.ok = (*bids[csm.B.Market])[csm.event.Data.Bids[csm.i][0]]; csm.ok {

			if csm.event.Data.Bids[csm.i][1] == (*bids[csm.B.Market])[csm.event.Data.Bids[csm.i][0]] {
				delete((*bids[csm.B.Market]), csm.event.Data.Bids[csm.i][0])
				(*csm.newBids)[csm.event.Data.Bids[csm.i][0]] = csm.event.Data.Bids[csm.i][1]
				continue
			}
			csm.B.Price = csm.event.Data.Bids[csm.i][0]
			delete((*bids[csm.B.Market]), csm.event.Data.Bids[csm.i][0])
			Del(csm, &csm.previousVolume)
		}

		csm.B.Price = csm.event.Data.Bids[csm.i][0]

		Upsert(csm, &csm.event.Data.Bids[csm.i][1])

		(*csm.newBids)[csm.event.Data.Bids[csm.i][0]] = csm.event.Data.Bids[csm.i][1]
	}

	// exchange the local with the global
	*asks[csm.B.Market], *csm.newAsks = *csm.newAsks, *asks[csm.B.Market]
	*bids[csm.B.Market], *csm.newBids = *csm.newBids, *bids[csm.B.Market]

	// // asks

	(*mutexMarkets)[csm.B.Market].Unlock()

	csm.B.Action = "1"
	// empty local asks and bids
	for csm.keyMap = range *csm.newAsks {
		csm.tmpString = (*csm.newAsks)[csm.keyMap]
		csm.B.Price = csm.keyMap

		Del(csm, &csm.tmpString)
		delete((*csm.newAsks), csm.keyMap)
	}
	csm.B.Action = "0"
	for csm.keyMap = range *csm.newBids {
		csm.tmpString = (*csm.newBids)[csm.keyMap]
		csm.B.Price = csm.keyMap

		Del(csm, &csm.tmpString)
		delete((*csm.newBids), csm.keyMap)
	}
}

func Del(csm *Consumer, volume *string) {
	// 000000
	csm.B.Volume = strings.Replace(*volume, ".", "", 1)
	disruptor.Del(&csm.B)

}
func Upsert(csm *Consumer, volume *string) {
	// 000000
	csm.B.Volume = strings.Replace(*volume, ".", "", 1)
	generateScore(csm)
	disruptor.Upsert(&csm.B)
}

func generateScore(csm *Consumer) error {
	csm.tmpFloat, csm.err = strconv.ParseFloat(csm.B.Price, 64)
	if csm.err != nil {
		return csm.err
	}

	csm.B.Score = strings.Replace(strconv.FormatFloat(csm.tmpFloat+(csm.tmpFloat*mapFees[csm.B.Market]), 'f', 2, 64), ".", "", 1)
	return nil
}
