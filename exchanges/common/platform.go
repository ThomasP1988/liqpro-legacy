package common

// Platform each connected platform should implement this interface
type Platform interface {
	Name() string
	Buy(symbol string, price float64, quantity float64) (*OrderResponse, error)
	Sell(symbol string, price float64, quantity float64) (*OrderResponse, error)
}

// OrderResponse format platform order answer with this struct
type OrderResponse struct {
	Price            float64
	QuantityAsked    float64
	QuantityExecuted float64
	Platform         string
	OrderID          string
}
