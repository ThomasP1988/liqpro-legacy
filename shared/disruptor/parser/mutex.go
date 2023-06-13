package socketredis

import "sync"

// SetMutexes we need a mutex per market
func SetMutexes(markets *[]string) *MutexMarkets {
	mutexes := MutexMarkets{}
	for i := 0; i < len(*markets); i++ {
		mutexes[(*markets)[i]] = &sync.Mutex{}
	}
	return &mutexes
}

// MutexMarkets  need mutexes for LMAX
type MutexMarkets = map[string]*sync.Mutex
