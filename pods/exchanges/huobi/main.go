package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"

	huobi "liqpro/exchanges/huobi"
	disruptor "liqpro/shared/disruptor/parser"

	"github.com/fasthttp/websocket"
)

var mapChannels map[string]string = map[string]string{
	"btcusdt": "btcusd",
	"ethusdt": "ethusd",
}
var markets = []string{
	"btcusdt",
	"ethusdt",
}
var orderbookSnapshot = map[string]bool{}
var lastSeq = map[string]int64{}

var lastProcessedEventMicroTimestamp int64 = 0
var mapFees map[string]*Fees = map[string]*Fees{}

var asks = map[string]map[string]string{}
var bids = map[string]map[string]string{}

var mutexMarkets *disruptor.MutexMarkets

var socketConnection *websocket.Conn

var startingMessageUpdate = []byte("{\"ch")
var startingMessageSnapshot = []byte("{\"id\":\"0")
var startingMessagePing = []byte("{\"ping")

//SocketWriteMutex prevent concurrent writing
var SocketWriteMutex = &sync.Mutex{}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	fees, err := huobi.GetFees(markets)
	println("fees", fees)
	if err != nil {
		println("huobi.GetFees", err)
		panic(err)
	}

	channels := []string{}

	for _, value := range mapChannels {
		channels = append(channels, value)
		asks[value] = map[string]string{}
		bids[value] = map[string]string{}
		lastSeq[value] = 0
		mapFees[value] = &Fees{
			Taker: 0,
			Maker: 0,
		}
	}

	for _, value := range fees.Data {
		if _, ok := mapFees[value.Symbol]; ok {
			mapFees[value.Symbol].Maker, _ = strconv.ParseFloat(value.ActualMakerRate, 64)
			mapFees[value.Symbol].Taker, _ = strconv.ParseFloat(value.ActualTakerRate, 64)
		}
	}

	mutexMarkets = disruptor.SetMutexes(&channels)

	go huobi.WsDepthServe(markets, huobi.WsDepthHandler(SocketIntercept))

	for {
		select {
		case <-interrupt:
			fmt.Println("interrupt")
			return
		}
	}
}

func parser(msg *[]byte, csm *Consumer) {
	csm.gr, csm.err = gzip.NewReader(bytes.NewBuffer(*msg))
	if csm.err != nil {
		fmt.Println("can't un-gzip", csm.err)
	}
	defer csm.gr.Close()
	// fmt.Println("csm.data", string(csm.data))
	csm.data, csm.err = ioutil.ReadAll(csm.gr)
	println("message", string(csm.data))
	if csm.err != nil {
		log.Println("error:", csm.err)
		return
	}

	if bytes.HasPrefix(csm.data, startingMessageUpdate) { // update

		csm.err = json.Unmarshal(csm.data, csm.event)
		if csm.err != nil {
			fmt.Println("error unmarshalling", string(csm.data), csm.err)
		} else {
			if csm.event.Channel != "" {
				csm.B.Market = mapChannels[strings.Split(csm.event.Channel, ".")[1]]
				ProcessEvent(csm)
			}
		}
	} else if bytes.HasPrefix(csm.data, startingMessageSnapshot) { // snapshot
		csm.err = json.Unmarshal(csm.data, csm.snap)
		if csm.err != nil {
			fmt.Println("error unmarshalling", string(csm.data), csm.err)
		} else {
			if csm.snap.Channel != "" {
				csm.B.Market = mapChannels[strings.Split(csm.snap.Channel, ".")[1]]
				ProcessSnapshot(csm)
			}
		}
	} else if bytes.HasPrefix(csm.data, startingMessagePing) {
		// send pong
		csm.err = json.Unmarshal(csm.data, csm.ping)
		if csm.err != nil {
			fmt.Println("error unmarshalling ping", string(csm.data), csm.err)
		}

		csm.byteJSONPong, csm.err = json.Marshal(csm.pong)
		SocketWriteMutex.Lock()
		csm.err = socketConnection.WriteMessage(websocket.TextMessage, csm.byteJSONPong)
		SocketWriteMutex.Unlock()
		if csm.err != nil {
			fmt.Println("error sending pong", csm.pong, csm.err)
		}
		return
	}

}

// SocketIntercept we need to access the socket to answer ping message
func SocketIntercept(c *websocket.Conn) {
	socketConnection = c
	wsDepthHandler := disruptor.Disruptor(NewConsumer())
	wsDepthHandler(c)
}

// ProcessSnapshot process event to populate map and DB
func ProcessSnapshot(csm *Consumer) {
	(*mutexMarkets)[csm.B.Market].Lock()
	defer (*mutexMarkets)[csm.B.Market].Unlock()
	csm.B.Action = "1"
	csm.dataLen = len(csm.snap.Data.Asks)
	fmt.Println("ask length", csm.dataLen)
	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {
		csm.B.Price = fmt.Sprintf("%.8f", csm.snap.Data.Asks[csm.i][0])
		if csm.previousVolume, csm.ok = asks[csm.B.Market][csm.B.Price]; csm.ok {
			csm.B.Volume = csm.previousVolume
			Del(csm)
		}

		csm.B.Volume = fmt.Sprintf("%.8f", csm.snap.Data.Asks[csm.i][1])
		if csm.snap.Data.Asks[csm.i][1] != 0 {
			Upsert(csm)
		}
		asks[csm.B.Market][csm.B.Price] = csm.B.Volume
	}
	csm.B.Action = "0"
	csm.dataLen = len(csm.snap.Data.Bids)

	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {

		csm.B.Price = fmt.Sprintf("%.8f", csm.snap.Data.Bids[csm.i][0])

		if csm.previousVolume, csm.ok = bids[csm.B.Market][csm.B.Price]; csm.ok {
			csm.B.Volume = csm.previousVolume
			Del(csm)
		}

		csm.B.Volume = fmt.Sprintf("%.8f", csm.snap.Data.Bids[csm.i][1])
		if csm.snap.Data.Bids[csm.i][1] != 0 {
			Upsert(csm)
		}
		bids[csm.B.Market][csm.B.Price] = csm.B.Volume
	}
}

// ProcessEvent process event to populate map and DB
func ProcessEvent(csm *Consumer) {
	(*mutexMarkets)[csm.B.Market].Lock()
	defer (*mutexMarkets)[csm.B.Market].Unlock()

	if lastSeq[csm.B.Market] != 0 && lastSeq[csm.B.Market] != csm.event.Tick.PreviousSequenceNumber {
		println("ON A RATE UN PUTAIN DE MESSAGE")
	}

	lastSeq[csm.B.Market] = csm.event.Tick.SequenceNumber

	csm.B.Action = "1"
	csm.dataLen = len(csm.event.Tick.Asks)

	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {

		csm.B.Price = fmt.Sprintf("%.8f", csm.event.Tick.Asks[csm.i][0])

		if csm.previousVolume, csm.ok = asks[csm.B.Market][csm.B.Price]; csm.ok {
			csm.B.Volume = csm.previousVolume
			Del(csm)
		}
		csm.B.Volume = fmt.Sprintf("%.8f", csm.event.Tick.Asks[csm.i][1])
		if csm.event.Tick.Asks[csm.i][1] != 0 {
			Upsert(csm)
			asks[csm.B.Market][csm.B.Price] = csm.B.Volume
		} else {
			test, test1 := asks[csm.B.Market][csm.B.Price]
			fmt.Print("DELETE", test, test1)
			delete(asks[csm.B.Market], csm.B.Price)
		}
	}

	csm.B.Action = "0"
	csm.dataLen = len(csm.event.Tick.Bids)

	for csm.i = 0; csm.i < csm.dataLen; csm.i++ {

		csm.B.Price = fmt.Sprintf("%.8f", csm.event.Tick.Bids[csm.i][0])

		if csm.previousVolume, csm.ok = bids[csm.B.Market][csm.B.Price]; csm.ok {
			csm.B.Volume = csm.previousVolume
			Del(csm)
		}
		csm.B.Volume = fmt.Sprintf("%.8f", csm.event.Tick.Bids[csm.i][1])
		if csm.event.Tick.Bids[csm.i][1] != 0 {
			Upsert(csm)
			bids[csm.B.Market][csm.B.Price] = csm.B.Volume
		} else {
			delete(bids[csm.B.Market], csm.B.Price)
		}
	}
	println("asks: ", len(asks[csm.B.Market]), "bids: ", len(bids[csm.B.Market]))

	csm.event.Tick.Bids = csm.event.Tick.Bids[:0]
	csm.event.Tick.Asks = csm.event.Tick.Asks[:0]
}

func Del(csm *Consumer) {
	// 000000
	csm.B.Volume = strings.Replace(csm.B.Volume, ".", "", 1)
	disruptor.Del(&csm.B)

}
func Upsert(csm *Consumer) {
	// 000000
	csm.B.Volume = strings.Replace(csm.B.Volume, ".", "", 1)
	generateScore(csm)
	disruptor.Upsert(&csm.B)
}

// TODO: optimised by reducing the high number of conversion
func generateScore(csm *Consumer) error {
	csm.tmpFloat, csm.err = strconv.ParseFloat(csm.B.Price, 64)
	if csm.err != nil {
		return csm.err
	}

	csm.B.Score = strings.Replace(strconv.FormatFloat(csm.tmpFloat+(csm.tmpFloat*mapFees[csm.B.Market].Taker), 'f', 2, 64), ".", "", 1)
	return nil
}
