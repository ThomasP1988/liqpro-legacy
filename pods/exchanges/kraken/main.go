package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"time"

	kraken "liqpro/exchanges/kraken"

	disruptor "liqpro/shared/disruptor/parser"
)

var platformName = "kraken"

var mapChannels map[string]string = map[string]string{
	"XBT/USD": "btcusd",
	"XBT/EUR": "btceur",
	"ETH/EUR": "btceur",
}
var markets = []string{"XBT/USD", "XBT/EUR", "ETH/EUR"}

var alphaMarkets []string = []string{}
var channels []string = []string{}

var (
	it   int
	clnr *disruptor.Cleaner
)

var lastProcessedEventMicroTimestamp int64 = 0

var asks = map[string]map[string]string{}
var bids = map[string]map[string]string{}

var mutexMarkets *disruptor.MutexMarkets

var startJsonBracket = []byte("{")
var startJsonSquareBracket = []byte("[")
var startJsonChannelId = []byte("{\"channelID")

func main() {
	defer CleanDB()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for their, our := range mapChannels {
		alphaMarkets = append(alphaMarkets, strings.Replace(their, "/", "", 1))
		channels = append(channels, our)
		asks[our] = map[string]string{}
		bids[our] = map[string]string{}
	}

	kraken.GetBalance()
	rslt, err := kraken.GetFees(markets)

	if err != nil {
		panic(err)
	}

	println("rslt", rslt)

	mutexMarkets = disruptor.SetMutexes(&channels)
	wsDepthHandler := disruptor.Disruptor(
		NewConsumer(),
		NewConsumer(),
	)

	var onConnected = func() {
		debug.SetGCPercent(-1)
	}

	go kraken.WsDepthServe(markets, kraken.WsDepthHandler(wsDepthHandler), &onConnected)

	for {
		select {
		case <-interrupt:
			fmt.Println("interrupt")
			return
		}
	}

}

func parser(msg *[]byte, csm *Consumer) {
	if bytes.HasPrefix(*msg, startJsonBracket) {
		if bytes.HasPrefix(*msg, startJsonChannelId) {
			csm.err = json.Unmarshal(*msg, csm.event)
			if csm.err != nil {
				fmt.Println("error unmarshalling", string(*msg), csm.err)
			} else {
				fmt.Println(csm.event)
			}
		}
	} else if bytes.HasPrefix(*msg, startJsonSquareBracket) {

		csm.err = json.Unmarshal(*msg, csm.data)
		if csm.err != nil {
			fmt.Println(csm.err)

		}

		_, csm.err = kraken.ParseData(csm.data.Data, csm.data.Pair, csm.dataParser)
		if csm.err != nil {
			return
		}
		csm.data.Data = csm.dataParser.Result
		csm.B.Market = mapChannels[csm.data.Pair]
		ParseEvent(csm)

	}
}

// ParseEvent populate our DB with this event
func ParseEvent(csm *Consumer) {
	(*mutexMarkets)[csm.B.Market].Lock()
	defer (*mutexMarkets)[csm.B.Market].Unlock()
	// fmt.Println("event", csm.data.Data.(kraken.OrderBookUpdate).Asks)
	csm.dataLen = len(csm.data.Data.(kraken.OrderBookUpdate).Asks)
	csm.B.Action = "1"
	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {

		csm.B.Price = csm.data.Data.(kraken.OrderBookUpdate).Asks[csm.i].Price

		if csm.previousVolume, csm.ok = asks[csm.B.Market][csm.B.Price]; csm.ok {
			csm.B.Volume = csm.previousVolume
			Del(csm)
		}
		if csm.data.Data.(kraken.OrderBookUpdate).Asks[csm.i].Volume != "0.00000000" {
			csm.B.Volume = csm.data.Data.(kraken.OrderBookUpdate).Asks[csm.i].Volume
			Upsert(csm)
			asks[csm.B.Market][csm.B.Price] = csm.B.Volume
		} else {
			delete(asks[csm.B.Market], csm.B.Price)
		}
	}
	csm.B.Action = "0"
	csm.dataLen = len(csm.data.Data.(kraken.OrderBookUpdate).Bids)
	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {

		csm.B.Price = csm.data.Data.(kraken.OrderBookUpdate).Bids[csm.i].Price

		if csm.previousVolume, csm.ok = bids[csm.B.Market][csm.B.Price]; csm.ok {
			csm.B.Volume = csm.previousVolume
			Del(csm)
		}
		if csm.data.Data.(kraken.OrderBookUpdate).Bids[csm.i].Volume != "0.00000000" {
			csm.B.Volume = csm.data.Data.(kraken.OrderBookUpdate).Bids[csm.i].Volume
			Upsert(csm)
			bids[csm.B.Market][csm.B.Price] = csm.B.Volume
		} else {
			delete(bids[csm.B.Market], csm.B.Price)
		}
	}
	// println("asks: ", len(asks[csm.B.Market]), "bids: ", len(bids[csm.B.Market]))

}

func Del(csm *Consumer) {
	// 000000

	disruptor.Del(&csm.B)

}
func Upsert(csm *Consumer) {
	// 000000

	disruptor.Upsert(&csm.B)
}

func CleanDB() {
	disruptor.StopDisruptor()
	time.Sleep(time.Second)
	clnr = disruptor.NewCleaner(platformName)

	for it = 0; it < len(channels); it++ {
		clnr.B.Market = channels[it]
		clnr.B.Action = "0"
		disruptor.CleanPlatform(clnr)
		clnr.B.Action = "1"
		disruptor.CleanPlatform(clnr)
	}
}
