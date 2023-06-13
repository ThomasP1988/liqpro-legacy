package order

// RequestArgs body request
type RequestArgs struct {
	Instrument string  `json:"instrument"`
	Quantity   float64 `json:"quantity"`
	Side       string  `json:"side"`
	ClientID   string  `json:"clientId,omitempty"`
}

// RequestResponse body request
type RequestResponse struct {
	Event        string  `json:"event"`
	Instrument   string  `json:"instrument"`
	Price        string  `json:"price"`
	PricePerUnit string  `json:"pricePerUnit"`
	Quantity     float64 `json:"quantity"`
	Side         string  `json:"side"`
	ClientID     string  `json:"clientId,omitempty"`
}
