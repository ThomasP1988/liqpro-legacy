package socketredis

import (
	"context"
	"strings"

	redis "github.com/go-redis/redis/v8"
)

const (
	batchSize    int64 = 10000
	securitySize int64 = 100 // some items can go up and down while looping, so better re-looping same items than letting wrong data
)

type Cleaner struct {
	rangeRedis *redis.ZRangeBy
	result     *redis.ZSliceCmd
	ctx        context.Context
	platform   string
	pair       string
	i          int
	dataLen    int
	entries    []redis.Z
	err        error
	B          ConsumerBase
	isEmpty    bool
	offset     int64
}

func CleanPlatform(clnr *Cleaner) {

	clnr.offset = 0

	clnr.rangeRedis.Min = "0"
	clnr.rangeRedis.Max = "+inf"
	clnr.rangeRedis.Count = batchSize

	for {
		clnr.rangeRedis.Offset = clnr.offset
		clnr.result = Pipeline.ZRangeByScoreWithScores(clnr.ctx, clnr.B.Market+":"+clnr.B.Action, clnr.rangeRedis)
		Pipeline.Exec(clnr.ctx)
		clnr.entries, clnr.err = clnr.result.Result()

		clnr.dataLen = len(clnr.entries)

		if clnr.dataLen == 0 {
			break
		}

		for clnr.i = 0; clnr.i < clnr.dataLen; clnr.i++ {
			keyStr := strings.Split(clnr.entries[clnr.i].Member.(string), ":")
			if keyStr[2] == clnr.B.Platform {
				clnr.B.Price = keyStr[1]
				clnr.B.Volume = keyStr[0]
				Del(&clnr.B)
				clnr.isEmpty = false
			}
		}

		clnr.offset += batchSize - securitySize
	}

}

func NewCleaner(platform string) *Cleaner {
	return &Cleaner{
		B: ConsumerBase{
			Platform: platform,
		},
		rangeRedis: &redis.ZRangeBy{},
		result:     &redis.ZSliceCmd{},
		ctx:        context.Background(),
	}
}
