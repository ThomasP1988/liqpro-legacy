package main

import "fmt"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients     map[*Client]bool
	channels    map[string]map[float64]map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	subscribe   chan Subscribe
	unsubscribe chan Subscribe
}

// NewHub centralise state
func NewHub() *Hub {
	return &Hub{
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		subscribe:   make(chan Subscribe),
		unsubscribe: make(chan Subscribe),
		clients:     make(map[*Client]bool),
		channels:    make(map[string]map[float64]map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			fmt.Println("register")
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				for _, sub := range client.subscriptions {
					h.UnsubscribeFc(&sub)
				}
				delete(h.clients, client)
				close(client.send)
			}
		case sub := <-h.subscribe:
			fmt.Println(sub)
			if _, ok := h.channels[sub.instrument]; !ok {
				h.channels[sub.instrument] = make(map[float64]map[*Client]bool)
				go TriggerInstrumentProcessor(sub.instrument)
			}
			if _, ok := h.channels[sub.instrument][sub.level]; !ok {
				h.channels[sub.instrument][sub.level] = make(map[*Client]bool)
			}
			h.channels[sub.instrument][sub.level][sub.client] = true
			sub.client.subscriptions = append(sub.client.subscriptions, sub)
		case sub := <-h.unsubscribe:
			h.UnsubscribeFc(&sub)
			// case message := <-h.broadcast:
			// 	for client := range h.clients {
			// 		select {
			// 		case client.send <- message:
			// 		default:
			// 			close(client.send)
			// 			delete(h.clients, client)
			// 		}
			// 	}
		}
	}
}

// UnsubscribeFc unsubscribe client to level or instrument
func (h *Hub) UnsubscribeFc(sub *Subscribe) {
	if _, ok := h.channels[sub.instrument][sub.level]; ok {
		delete(h.channels[sub.instrument][sub.level], sub.client)
		if len(h.channels[sub.instrument][sub.level]) == 0 {
			delete(h.channels[sub.instrument], sub.level)
		}
	}
}

// Subscribe struct to send to channels subscribe and unsubscribe
type Subscribe struct {
	client     *Client
	instrument string
	level      float64
}
