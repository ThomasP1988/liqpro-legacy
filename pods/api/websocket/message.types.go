package main

import "liqpro/shared/libs/order"

// MessageIn from client
type MessageIn struct {
	Event             string    `json:"event"`
	Instrument        string    `json:"instrument,omitempty"`
	Levels            []float64 `json:"levels,omitempty"`
	order.RequestArgs           // we need this struct to pass order
}
