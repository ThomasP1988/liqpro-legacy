package bitstamp

import (
	"encoding/json"
	"log"
)

// GetOrderBook Retrieve the order book for a specific market.
func GetOrderBook(symbol string) (*OrderBook, error) {
	body, err := Request("GET", "/order_book/"+symbol, map[string]string{"group": "1"}, "")
	if err != nil {
		log.Println("error:", err)
		return nil, err
	}
	ob := &OrderBook{}
	err = json.Unmarshal(body, ob)
	if err != nil {
		return nil, err
	}
	return ob, nil
}

// OrderBook bitstamp
type OrderBook struct {
	// Timestamp      int32      `json:"timestamp,string,omitempty"`
	Microtimestamp int64      `json:"microtimestamp,string"`
	Bids           [][]string `json:"bids,[][]string"`
	Asks           [][]string `json:"asks,[][]string"`
}
