package main

import (
	"encoding/json"
	"fmt"
	disruptor "liqpro/shared/disruptor/parser"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"

	"liqpro/exchanges/binance"
)

var mapChannels map[string]string = map[string]string{
	"btceur":  "btceur",
	"etheur":  "etheur",
	"btcusdt": "btcusd",
}
var markets = []string{"BTCEUR", "BTCUSDT", "ETHEUR"}

// var orderbookSnapshot = map[string]bool{}
// var orderbookSnapshotLastUpdateID = &map[string]int64{}

// var lastProcessedEventMicroTimestamp int64 = 0

var asks = map[string]*map[string]string{}
var bids = map[string]*map[string]string{}

var mutexMarkets *disruptor.MutexMarkets

// var eventQueue *EventQueue

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	account, err := binance.GetAccount()

	if err != nil {
		panic(err)
	}

	channels := []string{}

	for _, value := range mapChannels {
		channels = append(channels, value)
		asksMap := make(map[string]string, 30)
		bidsMap := make(map[string]string, 30)
		asks[value] = &asksMap
		bids[value] = &bidsMap
	}

	mutexMarkets = disruptor.SetMutexes(&channels)
	// eventQueue = SetEventQueue(&markets)

	wsDepthHandler := disruptor.Disruptor(
		NewConsumer(account),
		NewConsumer(account),
	)

	var onConnected = func() {
		debug.SetGCPercent(-1)
	}

	go binance.WsDepthServe(markets, binance.WsDepthHandler(wsDepthHandler), &onConnected)

	// for their, our := range mapChannels {
	// 	orderbookSnapshot[our] = false

	// 	orderbook, err := binance.GetOrderBook(strings.ToUpper(their))
	// 	if err != nil {
	// 		fmt.Println("Error getting orderbook:", err)
	// 		continue
	// 	}
	// 	(*orderbookSnapshotLastUpdateID)[our] = orderbook.LastUpdateID
	// 	for _, ask := range orderbook.Asks[:20] {
	// 		Upsert(our, "1", ask[0], ask[1])
	// 		(*asks[our])[ask[0]] = ask[1]
	// 	}
	// 	for _, bid := range orderbook.Bids[:20] {
	// 		Upsert(our, "0", bid[0], bid[1])
	// 		(*bids[our])[bid[0]] = bid[1]
	// 	}

	// }

	for {
		select {
		case <-interrupt:
			fmt.Println("interrupt")
			return
		}
	}

}

func parser(csm *Consumer) {
	csm.err = json.Unmarshal(*csm.message, csm.event)
	if csm.err != nil || csm.event.Stream == "" {
		fmt.Println("error unmarshalling", string(*csm.message), csm.err)
		return
	}

	csm.B.Market = mapChannels[strings.Split(csm.event.Stream, "@")[0]]
	(*mutexMarkets)[csm.B.Market].Lock()
	// if !orderbookSnapshot[csm.B.Market] {

	// 	//Push
	// 	lastUpdateID, ok := (*orderbookSnapshotLastUpdateID)[csm.B.Market]
	// 	if ok {
	// 		// lastUpdateID
	// 		for i := 0; i < len((*eventQueue)[csm.B.Market]); i++ {
	// 			if int64((*eventQueue)[csm.B.Market][i].Data.LastUpdateID) > lastUpdateID {
	// 				ProcessEvent(csm, false)
	// 			}
	// 		}
	// 		orderbookSnapshot[csm.B.Market] = true
	// 	} else {
	// 		(*eventQueue)[csm.B.Market] = append((*eventQueue)[csm.B.Market], csm.event)
	// 	}
	// 	(*mutexMarkets)[csm.B.Market].Unlock()

	// } else {
	// 	ProcessEvent(csm, true)
	// }
	ProcessEvent(csm, true)

}

// ProcessEvent process event to populate map and DB
// we will write data in the thread map and then interfere pointer, while cleaning the previous one
func ProcessEvent(csm *Consumer, unlock bool) {

	csm.dataLen = len(csm.event.Data.Asks)

	csm.B.Action = "1"
	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {
		if csm.previousVolume, csm.ok = (*asks[csm.B.Market])[csm.event.Data.Asks[csm.i][0]]; csm.ok {

			if csm.event.Data.Asks[csm.i][1] == (*asks[csm.B.Market])[csm.event.Data.Asks[csm.i][0]] {
				delete((*asks[csm.B.Market]), csm.event.Data.Asks[csm.i][0])
				(*csm.newAsks)[csm.event.Data.Asks[csm.i][0]] = csm.event.Data.Asks[csm.i][1]
				continue
			}

			delete((*asks[csm.B.Market]), csm.event.Data.Asks[csm.i][0])

			Del(csm, &csm.event.Data.Asks[csm.i][0], &csm.previousVolume)
		}

		Upsert(csm, &csm.event.Data.Asks[csm.i][0], &csm.event.Data.Asks[csm.i][1])
		(*csm.newAsks)[csm.event.Data.Asks[csm.i][0]] = csm.event.Data.Asks[csm.i][1]
	}

	csm.dataLen = len(csm.event.Data.Bids)
	csm.B.Action = "0"
	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {
		if csm.previousVolume, csm.ok = (*bids[csm.B.Market])[csm.event.Data.Bids[csm.i][0]]; csm.ok {

			if csm.event.Data.Bids[csm.i][1] == (*bids[csm.B.Market])[csm.event.Data.Bids[csm.i][0]] {
				delete((*bids[csm.B.Market]), csm.event.Data.Bids[csm.i][0])
				(*csm.newBids)[csm.event.Data.Bids[csm.i][0]] = csm.event.Data.Bids[csm.i][1]
				continue
			}

			delete((*bids[csm.B.Market]), csm.event.Data.Bids[csm.i][0])
			Del(csm, &csm.event.Data.Bids[csm.i][0], &csm.previousVolume)
		}

		Upsert(csm, &csm.event.Data.Bids[csm.i][0], &csm.event.Data.Bids[csm.i][1])

		(*csm.newBids)[csm.event.Data.Bids[csm.i][0]] = csm.event.Data.Bids[csm.i][1]
	}

	// exchange the local with the global
	*asks[csm.B.Market], *csm.newAsks = *csm.newAsks, *asks[csm.B.Market]
	*bids[csm.B.Market], *csm.newBids = *csm.newBids, *bids[csm.B.Market]

	// // asks

	if unlock {
		(*mutexMarkets)[csm.B.Market].Unlock()
	}
	csm.B.Action = "1"
	// empty local asks and bids
	for csm.keyMap = range *csm.newAsks {
		csm.tmpString = (*csm.newAsks)[csm.keyMap]
		Del(csm, &csm.keyMap, &csm.tmpString)
		delete((*csm.newAsks), csm.keyMap)
	}
	csm.B.Action = "0"
	for csm.keyMap = range *csm.newBids {
		csm.tmpString = (*csm.newBids)[csm.keyMap]
		Del(csm, &csm.keyMap, &csm.tmpString)
		delete((*csm.newBids), csm.keyMap)
	}

}

func Del(csm *Consumer, price *string, volume *string) {
	// 000000

	csm.B.Price = (*price)[:len(*price)-4]
	csm.B.Volume = strings.Replace(*volume, ".", "", 1)

	disruptor.Del(&csm.B)

}

func Upsert(csm *Consumer, price *string, volume *string) {
	// 000000

	csm.B.Price = (*price)[:len(*price)-4]
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
