package binance

import (
	"encoding/json"
	"log"
)

// GetOrderBook Retrieve the order book for a specific market.
func GetOrderBook(symbol string) (*OrderBook, error) {
	body, err := Request("GET", "/depth?symbol="+symbol+"&limit=20", map[string]string{"group": "1"}, "")

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

// OrderBook binance
type OrderBook struct {
	// Timestamp      int32      `json:"timestamp,string,omitempty"`
	LastUpdateID int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}
