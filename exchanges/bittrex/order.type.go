package bittrex

import "time"

// CreateOrderParams data type to place a new order
type CreateOrderParams struct {
	MarketSymbol string  `json:"marketSymbol"`
	Direction    string  `json:"direction"`
	Type         string  `json:"type"`
	Quantity     float64 `json:"quantity"`
	// Ceiling       float64 `json:"ceiling"`
	Limit         float64 `json:"limit"`
	TimeInForce   string  `json:"timeInForce"`
	ClientOrderID string  `json:"clientOrderId"`
	// UseAwards     string  `json:"useAwards"`
}

//OrderV3Response response from bittrex when passing an order
type OrderV3Response struct {
	ID           string  `json:"id"`
	MarketSymbol string  `json:"marketSymbol"`
	Direction    string  `json:"direction"`
	Type         string  `json:"type"`
	Quantity     float64 `json:"quantity"`
	Limit        float64 `json:"limit"`
	// Ceiling       float64   `json:"ceiling"`
	TimeInForce   string    `json:"timeInForce"`
	ClientOrderID string    `json:"clientOrderId"`
	FillQuantity  float64   `json:"fillQuantity"`
	Commission    float64   `json:"commission"`
	Proceeds      float64   `json:"proceeds"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ClosedAt      time.Time `json:"closedAt"`
}
