package bittrex

import (
	"encoding/json"
	"log"
)

// GetOrderBook Retrieve the order book for a specific market.
func GetOrderBook(symbol string) (*Orderbook, error) {
	body, err := Request("GET", "/markets/"+symbol+"/orderbook", map[string]string{"depth": "25"}, nil)
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}
	od := &Orderbook{}

	errUnmarshal := json.Unmarshal(*body, od)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return od, nil
}

// Orderbook orderbook snapshot we receive from REST API
type Orderbook struct {
	Bid []OrderbookData `json:"bid"`
	Ask []OrderbookData `json:"ask"`
}

// OrderbookData orderbook snapshot we receive from REST API
type OrderbookData struct {
	Quantity string `json:"quantity"`
	Rate     string `json:"rate"`
}
