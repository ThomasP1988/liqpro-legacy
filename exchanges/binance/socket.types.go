package binance

import (
	socket "github.com/fasthttp/websocket"
)

// WsDepthHandler handle websocket depth event
type WsDepthHandler func(c *socket.Conn)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsSubscribeMessage type for message to send to ws binance api
type WsSubscribeMessage struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int      `json:"id"`
}

// WsSubscribeMessageData type for message to send to ws binance api
type WsSubscribeMessageData struct {
	Channel string `json:"channel"`
}

// WsDepthEvent define websocket depth event
type WsDepthEvent struct {
	Stream string           `json:"stream"`
	Data   WsDepthEventData `json:"data"`
}

// WsDepthEventData define websocket depth event
type WsDepthEventData struct {
	Asks         [][]string `json:"asks"`
	Bids         [][]string `json:"bids"`
	LastUpdateID int        `json:"lastUpdateId"`
}

// WsDepthEvent define websocket depth event
// type WsDepthEvent struct {
// 	Event         string     `json:"e"`
// 	EventTime     int        `json:"E"`
// 	Symbol        string     `json:"s"`
// 	FirstUpdateID int        `json:"U"`
// 	FinalUpdateID int        `json:"s"`
// 	Bids          [][]string `json:"b"`
// 	Asks          [][]string `json:"a"`
// }
