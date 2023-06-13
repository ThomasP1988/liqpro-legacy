package main

import (
	"liqpro/exchanges/binance"
	disruptor "liqpro/shared/disruptor/parser"
)

// Consumer struct we have to implement to use the library
type Consumer struct {
	feeMaker       float64
	feeTaker       float64
	id             int64
	B              disruptor.ConsumerBase
	event          *binance.WsDepthEvent
	err            error
	previousVolume string
	newAsks        *map[string]string
	newBids        *map[string]string
	ok             bool
	keyMap         string
	i              int
	dataLen        int
	tmpString      string
	tmpFloat       float64
	message        *[]byte
}

// Variable to prevent two consumer executing the same message
var (
	consumerIdTotal int64 = 0
)

// Consume function to execute when thread is called
func (src Consumer) Consume(lower, upper int64) {
	for ; lower <= upper; lower++ {
		if lower%consumerIdTotal == src.id {
			src.message = disruptor.RingBuffer[lower&disruptor.BufferMask]
			// fmt.Println("message", message)
			parser(&src)
		}
	}
}

func NewConsumer(account *binance.Account) *Consumer {

	newAsks := make(map[string]string, 30)
	newBids := make(map[string]string, 30)

	csm := &Consumer{
		feeMaker: float64(account.MakerCommission) / 10000,
		feeTaker: float64(account.TakerCommission) / 10000,
		newAsks:  &newAsks,
		newBids:  &newBids,
		event:    &binance.WsDepthEvent{},
		B: disruptor.ConsumerBase{
			Platform: "binance",
		},
		id: consumerIdTotal,
	}

	consumerIdTotal++
	return csm

}
