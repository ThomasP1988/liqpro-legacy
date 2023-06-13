// +build !debug

package bittrex

import (
	"encoding/json"
	"fmt"
	common "liqpro/exchanges/common"
)

// Adapter exported PlatformAdapter
var Adapter common.Platform = PlatformAdapter{}

// PlatformAdapter implements common.Platform interface
type PlatformAdapter struct {
	common.Platform
}

// Name return const name of the platform
func (a PlatformAdapter) Name() string {
	return "bittrex"
}

// Buy buy on bittrex
func (a PlatformAdapter) Buy(symbol string, price float64, quantity float64) (*common.OrderResponse, error) {
	fmt.Println("RELEASE BUY")
	return a.order("BUY", symbol, price, quantity)
}

// Sell sell on bittrex
func (a PlatformAdapter) Sell(symbol string, price float64, quantity float64) (*common.OrderResponse, error) {
	fmt.Println("RELEASE SELL")
	return a.order("SELL", symbol, price, quantity)
}

func (a PlatformAdapter) order(side string, symbol string, price float64, quantity float64) (*common.OrderResponse, error) {

	payload := &CreateOrderParams{
		MarketSymbol: pairsOursToThem[symbol],
		Direction:    side,
		Type:         "LIMIT",
		Quantity:     quantity,
		Limit:        price,
		TimeInForce:  "IMMEDIATE_OR_CANCEL",
		// ClientOrderID: params.ClientOrderID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	response, errReponse := Request("POST", "/orders", map[string]string{}, &payloadBytes)

	if errReponse != nil {
		return nil, errReponse
	}
	fmt.Println("response", string(*response))
	result := &OrderV3Response{}

	err = json.Unmarshal(*response, result)

	if err != nil {
		return nil, err
	}

	return &common.OrderResponse{
		Price:            result.Limit,
		QuantityAsked:    quantity,
		QuantityExecuted: result.Quantity,
		Platform:         a.Name(),
		OrderID:          result.ID,
	}, nil

}
