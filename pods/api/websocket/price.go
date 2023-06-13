package main

import (
	"context"
	"fmt"
	config "liqpro/config"
	"math"
	"strconv"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
)

var (
	dbClient *redis.Client
	ctx      = context.Background()
	pipeline redis.Pipeliner
)

// ConnectToDB connect to redis
func ConnectToDB() {
	dbClient = redis.NewClient(&redis.Options{
		Addr:     "192.168.49.2:31928",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx = context.Background()
	pipeline = dbClient.Pipeline()
}

// SubscribePrice subscribe to range of price
func SubscribePrice(instrument string, cb SubscribeCB) {

	for {
		time.Sleep(time.Second / 2)
		resultAsks := pipeline.ZRangeByScoreWithScores(ctx, instrument+":1", &redis.ZRangeBy{
			Min:   "0",
			Max:   "+inf",
			Count: 100,
		})

		resultBids := pipeline.ZRevRangeByScoreWithScores(ctx, instrument+":0", &redis.ZRangeBy{
			Min:   "0",
			Max:   "+inf",
			Count: 100,
		})

		pipeline.Exec(ctx)
		cb(resultAsks, resultBids)
		// data, err := result.Result()

		if len((*TheHub).channels[instrument]) == 0 {
			delete((*TheHub).channels, instrument)
			break
		}
	}
}

type resultPrices = map[int]map[float64]float64

// CalculatePrices function to determine our prices
func CalculatePrices(marketResult *MarketResult, maxVolume float64, side int, results *resultPrices) error {
	fmt.Println("calculate prices")
	var levelIterator int = 0
	var totalAsksVolume float64 = 0
	var totalAsksPrice float64 = 0
	var totalBidsVolume float64 = 0
	var totalBidsPrice float64 = 0

	asks, errAsks := (*marketResult).Asks.Result()
	bids, errBids := (*marketResult).Bids.Result()

	if errAsks != nil || errBids != nil {
		fmt.Println("errAsks "+(*marketResult).market, errAsks)
		fmt.Println("errBids "+(*marketResult).market, errBids)
	}

	// fmt.Println(asks)
	// fmt.Println(bids)

	dataLen := len(asks)
	if side != 1 {
		for i := 0; i < dataLen; i++ {
			volumeStr := strings.Split(asks[i].Member.(string), ":")
			volume, err := strconv.ParseFloat(volumeStr[0], 64)
			if err != nil {
				return err
			}
			totalAsksVolume += volume
			totalAsksPrice += volume * asks[i].Score
			if totalAsksVolume > config.AuthorisedInstrumentsAndLevelsArray[marketResult.market][levelIterator] {
				// TODO: atm, we round to 2 decimal as we assume it is FIAT, in the future check if currency is FIAT or Crypto
				println("totalAsksVolume", totalAsksVolume)
				println("totalAsksPrice", totalAsksPrice)
				println("marketResult.market", marketResult.market)
				println("levelIterator", levelIterator)

				println("(*results)[0][config.AuthorisedInstrumentsAndLevelsArray[marketResult.market][levelIterator]]", (*results)[0][config.AuthorisedInstrumentsAndLevelsArray[marketResult.market][levelIterator]])
				(*results)[0][config.AuthorisedInstrumentsAndLevelsArray[marketResult.market][levelIterator]] = math.Ceil((totalAsksPrice/totalAsksVolume)*100) / 100
				levelIterator++
			}

			if maxVolume != 0 && maxVolume > totalAsksVolume {
				break
			}
		}
	}

	if side != 0 {
		levelIterator = 0
		dataLen = len(bids)

		for i := 0; i < dataLen; i++ {
			volumeStr := strings.Split(bids[i].Member.(string), ":")
			volume, err := strconv.ParseFloat(volumeStr[0], 64)
			if err != nil {
				return err
			}
			totalBidsVolume += volume
			totalBidsPrice += volume * bids[i].Score
			if totalBidsVolume > config.AuthorisedInstrumentsAndLevelsArray[marketResult.market][levelIterator] {
				// TODO: atm, we round to 2 decimal as we assume it is FIAT, in the future check if currency is FIAT or Crypto
				(*results)[1][config.AuthorisedInstrumentsAndLevelsArray[marketResult.market][levelIterator]] = math.Ceil((totalBidsPrice/totalBidsVolume)*100) / 100
				levelIterator++
			}

			if maxVolume != 0 && maxVolume > totalBidsVolume {
				break
			}
		}
	}
	return nil
}

// SubscribeCB subscribe callback when requesting prices
type SubscribeCB func(resultsAsks *redis.ZSliceCmd, resultBids *redis.ZSliceCmd)

// PriceFeed price structure we send to client
type PriceFeed struct {
	Event  string  `json:"event"`
	Ask    float64 `json:"ask"`
	Bid    float64 `json:"bid"`
	Market string  `json:"market"`
	Level  float64 `json:"level"`
}
