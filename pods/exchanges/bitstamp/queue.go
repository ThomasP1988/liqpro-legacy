package main

import (
	"liqpro/exchanges/bitstamp"
)

// TODO: factorise this function when generic will be available

// SetEventQueue instantiate the queue with our markets
func SetEventQueue(markets *[]string) *EventQueue {
	eventQueue := EventQueue{}
	for i := 0; i < len(*markets); i++ {
		eventQueue[(*markets)[i]] = []*bitstamp.WsDepthEvent{}
	}
	return &eventQueue
}

// EventQueue  need mutexes for LMAX
type EventQueue = map[string][]*bitstamp.WsDepthEvent
