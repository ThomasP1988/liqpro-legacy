package main

import (
	"encoding/json"
	"fmt"
	"sync"

	redis "github.com/go-redis/redis/v8"
	disruptor "github.com/smartystreets-prototypes/go-disruptor"
)

// variable for disruptor
const (
	BufferSize   = 32
	BufferMask   = BufferSize - 1
	Reservations = 1
)

var ringBuffer = [BufferSize]MarketResult{}
var myDisruptor = disruptor.New(
	disruptor.WithCapacity(BufferSize),
	disruptor.WithConsumerGroup(Consumer{}))

var MarketLastPrices map[string]*resultPrices = map[string]*resultPrices{}
var MarketLastPricesMutexes map[string]*sync.Mutex = map[string]*sync.Mutex{}

// TriggerInstrumentProcessor get the price from redis and add it to the disruptor
func TriggerInstrumentProcessor(instrument string) {
	var sequence int64 = 0

	SubscribePrice(instrument, func(resultsAsks *redis.ZSliceCmd, resultBids *redis.ZSliceCmd) {
		fmt.Println("callback")
		sequence = myDisruptor.Reserve(Reservations)

		ringBuffer[sequence&BufferMask] = MarketResult{
			market: instrument,
			Asks:   resultsAsks,
			Bids:   resultBids,
		}

		myDisruptor.Commit(sequence-Reservations+1, sequence)

	})

}

// Consumer struct we have to implement to use the library
type Consumer struct{}

// Consume function to execute when thread is called
func (src Consumer) Consume(lower, upper int64) {

	fmt.Println("ici")
	for ; lower <= upper; lower++ {
		marketResult := ringBuffer[lower&BufferMask]

		var results resultPrices = map[int]map[float64]float64{
			0: make(map[float64]float64),
			1: make(map[float64]float64),
		}

		err := CalculatePrices(&marketResult, 0, 2, &results)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("bids", results[0])
		fmt.Println("asks", results[1])
		// check prices changes
		var hasChanges bool = false

		// mutex
		if MarketLastPricesMutexes[marketResult.market] == nil {
			MarketLastPricesMutexes[marketResult.market] = &sync.Mutex{}
		}

		MarketLastPricesMutexes[marketResult.market].Lock()
		for level, clients := range (*TheHub).channels[marketResult.market] {
			if MarketLastPrices[marketResult.market] == nil || results[0][level] != (*MarketLastPrices[marketResult.market])[0][level] ||
				results[1][level] != (*MarketLastPrices[marketResult.market])[1][level] {

				feed := &PriceFeed{
					Event:  "price",
					Ask:    results[1][level],
					Bid:    results[0][level],
					Market: marketResult.market,
					Level:  level,
				}
				msg, err := json.Marshal(feed)

				if err != nil {
					fmt.Println("error marshalling pricefeed", err)
					continue
				}

				for client := range clients {
					fmt.Println("send")
					(*client).Send(msg)
				}
				hasChanges = true
			}
		}

		if hasChanges {
			MarketLastPrices[marketResult.market] = &results
		}
		MarketLastPricesMutexes[marketResult.market].Unlock()
		// end mutex
	}
}

// MarketResult format the result from redis to send it to the disruptor
type MarketResult struct {
	Asks   *redis.ZSliceCmd
	Bids   *redis.ZSliceCmd
	market string
}
