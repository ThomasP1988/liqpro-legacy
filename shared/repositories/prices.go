package repositories

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// Client to redis price variable
var (
	dbClient *redis.Client
	pipeline redis.Pipeliner
)

// ConnectRedisPrices connect to redis
func ConnectRedisPrices() {
	dbClient = redis.NewClient(&redis.Options{
		Addr:     "192.168.49.2:31928",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pipeline = dbClient.Pipeline()
}

// GetBestPricesBuy :)
func GetBestPricesBuy(pair string) ([]redis.Z, error) {
	CtxRedis := context.Background()

	var result *redis.ZSliceCmd
	result = pipeline.ZRangeByScoreWithScores(CtxRedis, pair+":1", &redis.ZRangeBy{
		Min:   "0",
		Max:   "+inf",
		Count: 100,
	})
	pipeline.Exec(CtxRedis)

	return result.Result()
}

// GetBestPricesSell :)
func GetBestPricesSell(pair string) ([]redis.Z, error) {
	CtxRedis := context.Background()

	var result *redis.ZSliceCmd
	result = pipeline.ZRevRangeByScoreWithScores(CtxRedis, pair+":0", &redis.ZRangeBy{
		Min:   "0",
		Max:   "+inf",
		Count: 100,
	})
	pipeline.Exec(CtxRedis)

	return result.Result()
}
