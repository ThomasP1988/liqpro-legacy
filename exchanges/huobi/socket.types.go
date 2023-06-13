package huobi

import (
	socket "github.com/fasthttp/websocket"
)

// WsDepthHandler handle websocket depth event
type WsDepthHandler func(c *socket.Conn)

// ErrHandler handles errors
type ErrHandler func(err error)

// Ping ping message
type Ping struct {
	Ping int `json:"ping"`
}

// Pong answer to ping message
type Pong struct {
	Pong string `json:"pong"`
}

// WsRequestMessage type for message to send to ws bitstamp api
type WsRequestMessage struct {
	Req string `json:"req"`
	ID  string `json:"id"`
}

// WsSubscribeMessage type for message to send to ws bitstamp api
type WsSubscribeMessage struct {
	Sub string `json:"sub"`
	ID  string `json:"id"`
}

// WsDepthEvent define websocket depth event
type WsDepthEvent struct {
	Channel   string           `json:"ch"`
	Timestamp int64            `json:"ts"`
	Tick      WsDepthEventTick `json:"tick"`
}

// WsDepthEventTick Tick websocket message
type WsDepthEventTick struct {
	SequenceNumber         int64       `json:"seqNum"`
	PreviousSequenceNumber int64       `json:"prevSeqNum,omitempty"`
	Bids                   [][]float64 `json:"bids,omitempty"`
	Asks                   [][]float64 `json:"asks,omitempty"`
}

// DepthDelta bid and ask differences from order book
type DepthDelta struct {
	Rate int64 `json:"rate"`
}

// WsSnapshotEvent define websocket depth event
type WsSnapshotEvent struct {
	Channel   string           `json:"rep"`
	Timestamp int64            `json:"ts"`
	Data      WsDepthEventTick `json:"data"`
}
