package main

import (
	"liqpro/exchanges/bittrex"
	disruptor "liqpro/shared/disruptor/parser"
)

// Consumer struct we have to implement to use the library
type Consumer struct {
	feeMaker       float64
	feeTaker       float64
	id             int64
	B              disruptor.ConsumerBase
	event          *bittrex.WsDepthEvent
	err            error
	previousVolume string
	microtimestamp int64
	ok             bool
	dataLen        int
	i              int
	tmpFloat       float64
	tmpString      string
	keyMap         string
	boolHelper     bool
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
			println("message", string(*message))
			parser(message, &src)
		}
	}
}

func NewConsumer(fee *bittrex.FeeRate) *Consumer {

	// the json structure requires those arrays
	bidDeltas := make([]bittrex.DepthDelta, 100)
	for i := range bidDeltas {
		bidDeltas[i] = bittrex.DepthDelta{}
	}
	askDeltas := make([]bittrex.DepthDelta, 100)
	for i := range askDeltas {
		askDeltas[i] = bittrex.DepthDelta{}
	}

	csm := &Consumer{
		feeMaker: float64(fee.Maker) / 10000,
		feeTaker: float64(fee.Taker) / 10000,
		event: &bittrex.WsDepthEvent{
			BidDeltas: bidDeltas,
			AskDeltas: askDeltas,
		},
		B: disruptor.ConsumerBase{
			Platform: "bittrex",
		},
		id: consumerIdTotal,
	}

	consumerIdTotal++
	return csm

}
