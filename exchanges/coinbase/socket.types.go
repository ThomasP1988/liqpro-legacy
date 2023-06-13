package coinbase

import (
	"encoding/json"

	socket "github.com/fasthttp/websocket"
)

// WsDepthHandler handle websocket depth event
type WsDepthHandler func(c *socket.Conn)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsSubscribeMessage type for message to send to ws bitstamp api
type WsSubscribeMessage struct {
	Type       string        `json:"type"`
	ProductIds []string      `json:"product_ids"`
	Channels   []interface{} `json:"channels"`
}

// WsDepthEvent define websocket depth event
type WsDepthEvent struct {
	Type      string     `json:"type"`
	ProductID string     `json:"product_id"`
	Changes   [][]string `json:"changes"` // [side, price, size]
	Time      string     `json:"time"`
}

// WsSnapShot define websocket snapshot
type WsSnapShot struct {
	Type      string     `json:"type"`
	ProductID string     `json:"product_id"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
}

// WsDepthEventData define websocket data depth event
type WsDepthEventData struct {
	Timestamp      int32           `json:"timestamp,string"`
	Microtimestamp int64           `json:"microtimestamp,string"`
	Bids           [][]json.Number `json:"bids,[][]string"`
	Asks           [][]json.Number `json:"asks,[][]string"`
}

// DepthDelta bid and ask differences from order book
type DepthDelta struct {
	Rate int64 `json:"rate"`
}
