package main

import (
	"liqpro/exchanges/kraken"
	disruptor "liqpro/shared/disruptor/parser"
)

// Consumer struct we have to implement to use the library
type Consumer struct {
	id             int64
	B              disruptor.ConsumerBase
	event          *kraken.WsDepthEvent
	data           *kraken.DataUpdate
	dataParser     *kraken.DataParser
	message        *[]byte
	err            error
	previousVolume string
	microtimestamp int64
	ok             bool
	dataLen        int
	i              int
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
			parser(src.message, &src)
		}
	}
}

func NewConsumer() *Consumer {

	csm := &Consumer{
		data:  &kraken.DataUpdate{},
		event: &kraken.WsDepthEvent{},
		B: disruptor.ConsumerBase{
			Platform: "kraken",
		},
		id:         consumerIdTotal,
		dataParser: &kraken.DataParser{},
	}

	consumerIdTotal++
	return csm

}
