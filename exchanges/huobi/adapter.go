// +build !debug

package huobi

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
	return "huobi"
}

// Buy buy on bitstamp
func (a PlatformAdapter) Buy(symbol string, price float64, quantity float64) (*common.OrderResponse, error) {
	return a.order("buy-ioc", pairsOursToThem[symbol], price, quantity)
}

// Sell sell on bitstamp
func (a PlatformAdapter) Sell(symbol string, price float64, quantity float64) (*common.OrderResponse, error) {
	return a.order("sell-ioc", pairsOursToThem[symbol], price, quantity)
}

func (a PlatformAdapter) order(side string, symbol string, price float64, quantity float64) (*common.OrderResponse, error) {

	order := &PlaceOrderRequest{
		AccountID: AccountID,
		Type:      side,
		Source:    "spot-api",
		Symbol:    "btcusdt",
		Price:     fmt.Sprintf("%.15f", price),
		// Amount:    fmt.Sprintf("%.15f",quantity),
		Amount: "0",
	}

	orderBytes, err := json.Marshal(order)
	if err != nil {
		fmt.Println("huobi order json marshal", err)
		return nil, err
	}

	response, errRequest := RequestSigned("POST", "/v1/order/orders/place", orderBytes, nil)
	fmt.Println("response", string(response))
	if errRequest != nil {
		fmt.Println("huobi order request error", errRequest)
		return nil, errRequest
	}

	// result := PlaceOrdersResponse{}
	// jsonErr = json.Unmarshal(response, &result)
	// if jsonErr != nil {
	// 	return nil, jsonErr
	// }

	return &common.OrderResponse{}, nil
}
