package main

import (
	"encoding/json"
	"fmt"
	"liqpro/exchanges/bittrex"
	disruptor "liqpro/shared/disruptor/parser"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
)

var mapChannels map[string]string = map[string]string{
	"ETH-EUR": "etheur",
	"BTC-EUR": "btceur",
	"BTC-USD": "btcusd",
}

var markets = []string{"BTC-EUR", "BTC-USD", "ETH-EUR"}
var orderbookSnapshot = map[string]bool{}
var processQueue = &map[string]bool{}

var asks = map[string]map[string]string{}
var bids = map[string]map[string]string{}

var mutexMarkets *disruptor.MutexMarkets
var eventQueue *EventQueue

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	fee, err := bittrex.GetFees()

	if err != nil {
		println("bittrex.GetFees()", err)
		panic(err)
	}

	channels := []string{}

	for _, value := range mapChannels {
		asks[value] = map[string]string{}
		bids[value] = map[string]string{}
		channels = append(channels, value)
	}

	mutexMarkets = disruptor.SetMutexes(&channels)
	eventQueue = SetEventQueue(&markets)
	wsDepthHandler := disruptor.DisruptorBytes(NewConsumer(fee))

	go bittrex.WsDepthServe(markets, bittrex.WsDepthHandler(wsDepthHandler))
	for their, our := range mapChannels {
		orderbookSnapshot[our] = false

		orderbook, err := bittrex.GetOrderBook(their)

		if err != nil {
			fmt.Println("Error getting orderbook:", err)
			continue
		}

		temporaryConsumer := &Consumer{
			B: disruptor.ConsumerBase{
				Platform: "bittrex",
				Market:   our,
			},
		}

		temporaryConsumer.B.Action = "1"

		for _, ask := range orderbook.Ask {
			Upsert(temporaryConsumer, &ask.Rate, &ask.Quantity)
			asks[our][ask.Rate] = ask.Quantity
		}

		temporaryConsumer.B.Action = "0"

		for _, bid := range orderbook.Bid {
			Upsert(temporaryConsumer, &bid.Rate, &bid.Quantity)
			bids[our][bid.Rate] = bid.Quantity
		}
		(*processQueue)[our] = true

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

	csm.err = json.Unmarshal(*msg, csm.event)
	if csm.err != nil {
		fmt.Println("error unmarshalling", string(*msg), csm.err)
	}
	csm.B.Market = mapChannels[csm.event.MarketSymbol]
	(*mutexMarkets)[csm.B.Market].Lock()
	defer (*mutexMarkets)[csm.B.Market].Unlock()
	if !orderbookSnapshot[csm.B.Market] {
		//Push
		_, csm.ok = (*processQueue)[csm.B.Market]
		if csm.ok {
			// microtimestamp

			for csm.i = 0; csm.i < len((*eventQueue)[csm.B.Market]); csm.i++ {

				ProcessEvent(csm, (*eventQueue)[csm.B.Market][csm.i])

			}
			orderbookSnapshot[csm.B.Market] = true

			csm.boolHelper = true
			for _, value := range orderbookSnapshot {
				if !value {
					csm.boolHelper = false
				}
			}

			if csm.boolHelper {
				debug.SetGCPercent(-1)
			}

		} else {
			event := &bittrex.WsDepthEvent{}
			csm.err = json.Unmarshal(*msg, event)
			if csm.err != nil {
				fmt.Println("error unmarshalling", string(*msg), csm.err)
			}
			(*eventQueue)[csm.B.Market] = append((*eventQueue)[csm.B.Market], event)
		}

	} else {
		ProcessEvent(csm, csm.event)
	}
}

// ProcessEvent process event to populate map and DB
func ProcessEvent(csm *Consumer, event *bittrex.WsDepthEvent) {

	csm.B.Action = "1"
	csm.dataLen = len(event.AskDeltas)
	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {
		if csm.previousVolume, csm.ok = asks[csm.B.Market][event.AskDeltas[csm.i].Rate]; csm.ok {
			Del(csm, &event.AskDeltas[csm.i].Rate, &csm.B.Volume)
		}
		if event.AskDeltas[csm.i].Quantity != "0" {
			Upsert(csm, &event.AskDeltas[csm.i].Rate, &event.AskDeltas[csm.i].Quantity)
			asks[csm.B.Market][event.AskDeltas[csm.i].Rate] = event.AskDeltas[csm.i].Quantity
		} else {
			delete(asks[csm.B.Market], event.AskDeltas[csm.i].Rate)
		}
	}

	csm.dataLen = len(event.BidDeltas)
	csm.B.Action = "0"
	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {
		if csm.previousVolume, csm.ok = bids[csm.B.Market][event.BidDeltas[csm.i].Rate]; csm.ok {
			Del(csm, &event.BidDeltas[csm.i].Rate, &csm.previousVolume)
		}
		if event.BidDeltas[csm.i].Quantity != "0" {
			Upsert(csm, &event.BidDeltas[csm.i].Rate, &event.BidDeltas[csm.i].Quantity)
			bids[csm.B.Market][event.BidDeltas[csm.i].Rate] = event.BidDeltas[csm.i].Quantity
		} else {
			delete(bids[csm.B.Market], event.BidDeltas[csm.i].Rate)
		}

	}

	event.AskDeltas = event.AskDeltas[:0]
	event.BidDeltas = event.BidDeltas[:0]

}

func Del(csm *Consumer, price *string, volume *string) {
	// 000000

	csm.B.Price = (*price)[:len(*price)-6]
	csm.B.Volume = strings.Replace(*volume, ".", "", 1)

	disruptor.Del(&csm.B)

}
func Upsert(csm *Consumer, price *string, volume *string) {
	// 000000
	csm.B.Price = (*price)[:len(*price)-6]
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
