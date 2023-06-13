package bittrex

import "log"

// ListMarkets List markets.
func ListMarkets() {
	body, err := Request("GET", "/markets", nil, nil)
	if err != nil {
		log.Println("error:", err)
		return
	}
	log.Printf("result: %s", string(*body))
}
