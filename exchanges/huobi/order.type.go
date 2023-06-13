package huobi

// PlaceOrderRequest data to send to place order
type PlaceOrderRequest struct {
	AccountID     string `json:"account-id"`
	Symbol        string `json:"symbol"`
	Type          string `json:"type"`
	Amount        string `json:"amount"`
	Price         string `json:"price,omitempty"`
	Source        string `json:"source,omitempty"`
	ClientOrderID string `json:"client-order-id,omitempty"`
	StopPrice     string `json:"stop-price,omitempty"`
	Operator      string `json:"operator,omitempty"`
}

// PlaceOrdersResponse response from Huobi
type PlaceOrdersResponse struct {
	Status       string `json:"status"`
	ErrorCode    string `json:"err-code"`
	ErrorMessage string `json:"err-msg"`
	Data         []PlaceOrderResult
}

// PlaceOrderResult sub interface from response
type PlaceOrderResult struct {
	OrderID       int64  `json:"order-id"`
	ClientOrderID string `json:"client-order-id"`
	ErrorCode     string `json:"err-code"`
	ErrorMessage  string `json:"err-msg"`
}
