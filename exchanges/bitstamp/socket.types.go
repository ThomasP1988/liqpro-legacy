package bitstamp

import (
	socket "github.com/fasthttp/websocket"
)

// WsDepthHandler handle websocket depth event
type WsDepthHandler func(c *socket.Conn)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsSubscribeMessage type for message to send to ws bitstamp api
type WsSubscribeMessage struct {
	Event string                 `json:"event"`
	Data  WsSubscribeMessageData `json:"data"`
}

// WsSubscribeMessageData type for message to send to ws bitstamp api
type WsSubscribeMessageData struct {
	Channel string `json:"channel"`
}

// WsDepthEvent define websocket depth event
type WsDepthEvent struct {
	Data    WsDepthEventData `json:"data"`
	Channel string           `json:"channel"`
	Event   string           `json:"event"`
}

// WsDepthEventData define websocket data depth event
type WsDepthEventData struct {
	Timestamp      int32      `json:"timestamp,string"`
	Microtimestamp int64      `json:"microtimestamp,string"`
	Bids           [][]string `json:"bids,[][]string"`
	Asks           [][]string `json:"asks,[][]string"`
}

// DepthDelta bid and ask differences from order book
type DepthDelta struct {
	Rate int64 `json:"rate"`
}
