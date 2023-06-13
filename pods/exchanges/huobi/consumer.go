package main

import (
	"compress/gzip"
	"liqpro/exchanges/huobi"
	disruptor "liqpro/shared/disruptor/parser"
)

// Consumer struct we have to implement to use the library
type Consumer struct {
	id             int64
	B              disruptor.ConsumerBase
	event          *huobi.WsDepthEvent
	snap           *huobi.WsSnapshotEvent
	ping           *huobi.Ping
	pong           *huobi.Pong
	byteJSONPong   []byte
	gr             *gzip.Reader
	dataLen        int
	data           []byte
	err            error
	previousVolume string
	tmpFloat       float64
	ok             bool
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
			message := disruptor.RingBuffer[lower&disruptor.BufferMask]
			parser(message, &src)
		}
	}
}

func NewConsumer() *Consumer {

	// the json structure requires those arrays
	changes := make([][]string, 100)
	for i := range changes {
		changes[i] = make([]string, 2)
	}

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
		snap:  &huobi.WsSnapshotEvent{},
		event: &huobi.WsDepthEvent{},
		ping:  &huobi.Ping{},
		pong:  &huobi.Pong{},
		B: disruptor.ConsumerBase{
			Platform: "huobi",
		},
		id: consumerIdTotal,
	}

	consumerIdTotal++
	return csm

}
