package main

import (
	"liqpro/exchanges/coinbase"
	disruptor "liqpro/shared/disruptor/parser"
	"strconv"
)

// Consumer struct we have to implement to use the library
type Consumer struct {
	feeMaker       float64
	feeTaker       float64
	id             int64
	B              disruptor.ConsumerBase
	event          *coinbase.WsDepthEvent
	snap           *coinbase.WsSnapShot
	err            error
	previousVolume string
	microtimestamp int64
	ok             bool
	dataLen        int
	i              int
	tmpString      string
	tmpFloat       float64
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
			parser(message, &src)
		}
	}
}

func NewConsumer(fee *coinbase.Fees) *Consumer {

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

	maker, _ := strconv.ParseFloat(fee.Maker, 64)
	taker, _ := strconv.ParseFloat(fee.Taker, 64)

	csm := &Consumer{
		feeMaker: maker,
		feeTaker: taker,
		snap: &coinbase.WsSnapShot{
			Bids: bids,
			Asks: asks,
		},
		event: &coinbase.WsDepthEvent{
			Changes: changes,
		},
		B: disruptor.ConsumerBase{
			Platform: "coinbase",
		},
		id: consumerIdTotal,
	}

	consumerIdTotal++
	return csm

}
