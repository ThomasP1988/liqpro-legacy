// +build !debug

package bitstamp

import (
	"encoding/json"
	"fmt"
	common "liqpro/exchanges/common"
	"net/url"
)

// Adapter exported PlatformAdapter
var Adapter common.Platform = PlatformAdapter{}

// PlatformAdapter implements common.Platform interface
type PlatformAdapter struct {
	common.Platform
}

// Name return const name of the platform
func (a PlatformAdapter) Name() string {
	return "bitstamp"
}

// Buy buy on bitstamp
func (a PlatformAdapter) Buy(symbol string, price float64, quantity float64) (*common.OrderResponse, error) {
	return a.order("buy", symbol, price, quantity)
}

// Sell sell on bitstamp
func (a PlatformAdapter) Sell(symbol string, price float64, quantity float64) (*common.OrderResponse, error) {
	return a.order("sell", symbol, price, quantity)
}

func (a PlatformAdapter) order(side string, symbol string, price float64, quantity float64) (*common.OrderResponse, error) {
	// https://www.bitstamp.net/api/v2/buy/{currency_pair}/
	params := url.Values{}
	params.Set("price", fmt.Sprintf("%.15f", price))
	params.Set("amount", fmt.Sprintf("%.15f", quantity))
	params.Set("ioc_order", "True")

	response, err := RequestSigned("POST", "/"+side+"/"+symbol+"/", params)
	fmt.Println("response", string(response))
	if err != nil {
		return nil, err
	}

	result := &LimitOrderResult{}
	err = json.Unmarshal(response, result)

	if err != nil {
		return nil, err
	}

	return &common.OrderResponse{
		Price:            result.Price,
		QuantityAsked:    quantity,
		QuantityExecuted: result.Amount,
		Platform:         a.Name(),
		OrderID:          fmt.Sprint(result.Id),
	}, nil
}
