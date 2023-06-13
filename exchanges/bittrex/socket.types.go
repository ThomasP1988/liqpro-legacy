package bittrex

// WsDepthHandler handle websocket depth event
type WsDepthHandler func(msg *[]byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsDepthEvent define websocket depth event
type WsDepthEvent struct {
	MarketSymbol string       `json:"marketSymbol"`
	Depth        int          `json:"depth"`
	Sequence     int          `json:"sequence"`
	BidDeltas    []DepthDelta `json:"bidDeltas"`
	AskDeltas    []DepthDelta `json:"askDeltas"`
}

// DepthDelta bid and ask differences from order book
type DepthDelta struct {
	Quantity string `json:"quantity"`
	Rate     string `json:"rate"`
}
