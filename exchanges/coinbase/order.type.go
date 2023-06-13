package coinbase

import "time"

// Order type to pass order to API
type Order struct {
	Type      string `json:"type"`
	Size      string `json:"size,omitempty"`
	Side      string `json:"side"`
	ProductID string `json:"product_id"`
	ClientOID string `json:"client_oid,omitempty"`
	Stp       string `json:"stp,omitempty"`
	Stop      string `json:"stop,omitempty"`
	StopPrice string `json:"stop_price,omitempty"`
	// Limit Order
	Price       string `json:"price,omitempty"`
	TimeInForce string `json:"time_in_force,omitempty"`
	PostOnly    bool   `json:"post_only,omitempty"`
	CancelAfter string `json:"cancel_after,omitempty"`
	// Market Order
	Funds          string `json:"funds,omitempty"`
	SpecifiedFunds string `json:"specified_funds,omitempty"`
}

// OrderResponse from API
type OrderResponse struct {
	ID            string    `json:"id"`
	Status        string    `json:"status,omitempty"`
	Settled       bool      `json:"settled,omitempty"`
	DoneReason    string    `json:"done_reason,omitempty"`
	DoneAt        time.Time `json:"done_at,string,omitempty"`
	CreatedAt     time.Time `json:"created_at,string,omitempty"`
	FillFees      string    `json:"fill_fees,omitempty"`
	FilledSize    string    `json:"filled_size,omitempty"`
	ExecutedValue string    `json:"executed_value,omitempty"`
}
