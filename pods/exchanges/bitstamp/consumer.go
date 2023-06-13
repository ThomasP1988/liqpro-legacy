package main

import (
	"liqpro/exchanges/bitstamp"
	disruptor "liqpro/shared/disruptor/parser"
)

// Consumer struct we have to implement to use the library
type Consumer struct {
	id             int64
	B              disruptor.ConsumerBase
	event          *bitstamp.WsDepthEvent
	newAsks        *map[string]string
	newBids        *map[string]string
	err            error
	previousVolume string
	microtimestamp int64
	ok             bool
	dataLen        int
	i              int
	tmpString      string
	tmpFloat       float64
	keyMap         string
}

// Variable to prevent two consumer executing the same message
var (
	consumerIdTotal int64 = 0
)

// Consume function to execute when thread is called
func (src Consumer) Consume(lower, upper int64) {
	for ; lower <= upper; lower++ {
		if lower%consumerIdTotal == src.id {
			message := disruptor.RingBuffer[lower&disruptor.BufferMask]
			parser(message, &src)
		}
	}
}

func NewConsumer() *Consumer {
	// local
	newAsks := make(map[string]string, 30)
	newBids := make(map[string]string, 30)

	// the json structure requires those arrays
	bids := make([][]string, 100)
	for i := range bids {
		bids[i] = make([]string, 2)
	}
	asks := make([][]string, 100)
	for i := range asks {
		asks[i] = make([]string, 2)
	}

	csm := &Consumer{
		newAsks: &newAsks,
		newBids: &newBids,
		event: &bitstamp.WsDepthEvent{
			Data: bitstamp.WsDepthEventData{
				Bids: bids,
				Asks: asks,
			},
		},
		B: disruptor.ConsumerBase{
			Platform: "bitstamp",
		},
		id: consumerIdTotal,
	}

	consumerIdTotal++
	return csm

}
