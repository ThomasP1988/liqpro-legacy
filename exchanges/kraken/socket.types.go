package kraken

import (
	socket "github.com/fasthttp/websocket"
)

// WsDepthHandler handle websocket depth event
type WsDepthHandler func(c *socket.Conn)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsSubscribeMessage type for message to send to ws bitstamp api
type WsSubscribeMessage struct {
	Event        string                         `json:"event"`
	Pair         []string                       `json:"pair"`
	Subscription WsSubscribeMessageSubscription `json:"subscription"`
}

// WsSubscribeMessageSubscription what kind of subscription you want
type WsSubscribeMessageSubscription struct {
	Depth int16 `json:"depth,omitempty"`
	// Interval int16  `json:"interval,omitempty"`
	Name string `json:"name,omitempty"`
	// Snapshot bool   `json:"snapshot,omitempty"`
}

// WsDepthEvent define websocket depth event
type WsDepthEvent struct {
	Type      string     `json:"type"`
	ProductID string     `json:"product_id"`
	Changes   [][]string `json:"changes"`
	Time      string     `json:"time"`
}

// WsSnapShot define websocket snapshot
type WsSnapShot struct {
	Type      string     `json:"type"`
	ProductID string     `json:"product_id"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
}

// DataUpdate - data structure of default Kraken WS update
type DataUpdate struct {
	ChannelID   int64
	Data        interface{}
	ChannelName string
	Pair        string
	Sequence    int64
}

// OrderBookUpdate - data structure for order book update
type OrderBookUpdate struct {
	Asks       []OrderBookItem
	Bids       []OrderBookItem
	IsSnapshot bool
	Pair       string
	CheckSum   string
}

// OrderBookItem - data structure for order book item
type OrderBookItem struct {
	Price     string
	Volume    string
	Time      float64
	Republish bool
}
