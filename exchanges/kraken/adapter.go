// +build !debug

package kraken

import (
	"fmt"
	common "liqpro/exchanges/common"
	"net/url"
	"strconv"
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

	params := &url.Values{
		"pair":      {pairsOursToThem[symbol]},
		"volume":    {strconv.FormatFloat(quantity, 'f', 8, 64)},
		"type":      {side},
		"ordertype": {"limit"},
	}

	response, errRequest := RequestSigned("AddOrder", params)
	fmt.Println("response", string(response))
	if errRequest != nil {
		fmt.Println("kraken order request error", errRequest)
		return nil, errRequest
	}
	// result := PlaceOrdersResponse{}
	// jsonErr = json.Unmarshal(response, &result)
	// if jsonErr != nil {
	// 	return nil, jsonErr
	// }

	return &common.OrderResponse{}, nil
}
