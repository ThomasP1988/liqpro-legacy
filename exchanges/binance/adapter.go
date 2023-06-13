// +build !debug

package binance

import (
	"encoding/json"
	"fmt"
	common "liqpro/exchanges/common"
	"strconv"
	"time"
)

// Adapter exported PlatformAdapter
var Adapter common.Platform = PlatformAdapter{}

// PlatformAdapter implements common.Platform interface
type PlatformAdapter struct {
	common.Platform
}

// Name return const name of the platform
func (a PlatformAdapter) Name() string {
	return "binance"
}

// Buy buy on binance
func (a PlatformAdapter) Buy(symbol string, price float64, quantity float64) (*common.OrderResponse, error) {
	fmt.Println("RELEASE BUY")
	return a.order("BUY", symbol, price, quantity)
}

// Sell sell on binance
func (a PlatformAdapter) Sell(symbol string, price float64, quantity float64) (*common.OrderResponse, error) {
	fmt.Println("RELEASE SELL")
	return a.order("SELL", symbol, price, quantity)
}

func (a PlatformAdapter) order(side string, symbol string, price float64, quantity float64) (*common.OrderResponse, error) {

	params := &map[string]string{
		"symbol":           pairsOursToThem[symbol],
		"side":             side,
		"price":            fmt.Sprintf("%.15f", price),
		"quantity":         fmt.Sprintf("%.15f", quantity),
		"type":             "LIMIT",
		"timestamp":        fmt.Sprint(time.Now().Unix() * 1000),
		"timeInForce":      "IOC", //Immediate Or Cancel: An order will try to fill the order as much as it can before the order expires.
		"newOrderRespType": "RESULT",
		"recvWindow":       "2000",
	}
	fmt.Println("params", params)
	response, err := RequestSigned("POST", "/order/test", params, "")
	fmt.Println("response", string(response))

	if err != nil {
		return nil, err
	}

	orderResponse := &OrderResponse{}

	err = json.Unmarshal(response, orderResponse)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	price, errPrice := strconv.ParseFloat(orderResponse.Price, 64)
	executedQty, errQty := strconv.ParseFloat(orderResponse.ExecutedQuantity, 64)
	if errPrice != nil && errQty != nil {
		fmt.Println("errPrice", errPrice)
		fmt.Println("errQty", errQty)
	}

	return &common.OrderResponse{
		Price:            price,
		QuantityAsked:    quantity,
		QuantityExecuted: executedQty,
		Platform:         a.Name(),
		OrderID:          orderResponse.ClientOrderID,
	}, nil
}
