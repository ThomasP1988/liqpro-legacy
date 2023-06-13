package binance

// OrderResponse define create order response
type OrderResponse struct {
	Symbol                   string `json:"symbol"`
	OrderID                  int64  `json:"orderId"`
	ClientOrderID            string `json:"clientOrderId"`
	TransactTime             int64  `json:"transactTime"`
	Price                    string `json:"price"`
	OrigQuantity             string `json:"origQty"`
	ExecutedQuantity         string `json:"executedQty"`
	CummulativeQuoteQuantity string `json:"cummulativeQuoteQty"`
	IsIsolated               bool   `json:"isIsolated"` // for isolated margin

	Status      string `json:"status"`
	TimeInForce string `json:"timeInForce"`
	Type        string `json:"type"`
	Side        string `json:"side"`

	// for order response is set to FULL
	Fills                 []*Fill `json:"fills"`
	MarginBuyBorrowAmount string  `json:"marginBuyBorrowAmount"` // for margin
	MarginBuyBorrowAsset  string  `json:"marginBuyBorrowAsset"`
}

// Fill may be returned in an array of fills in a CreateOrderResponse.
type Fill struct {
	Price           string `json:"price"`
	Quantity        string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
}
