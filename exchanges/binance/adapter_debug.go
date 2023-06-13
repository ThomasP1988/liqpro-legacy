// +build debug

// this debug version is here to mock the answer from binance
package binance

import (
	"crypto/rand"
	"encoding/base64"
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
	fmt.Println("DEBUG BUY")
	return a.order("BUY", symbol, price, quantity)
}

// Sell sell on binance
func (a PlatformAdapter) Sell(symbol string, price float64, quantity float64) (*common.OrderResponse, error) {
	fmt.Println("DEBUG SELL")
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
	}
	fmt.Println("params", params)
	response, err := RequestSigned("POST", "/order/test", params, "")
	fmt.Println("response", string(response))

	if err != nil {
		return nil, err
	}

	// DEBUG Mocking for debug purposes
	var executedQtyMock string
	// we want to try when its fully filled or partially
	if HasFilled {
		executedQtyMock = fmt.Sprintf("%.15f", quantity)
	} else {
		executedQtyMock = fmt.Sprintf("%.15f", quantity/2)
	}
	fmt.Println("executedQtyMock", executedQtyMock)
	HasFilled = !HasFilled

	orderResponse := &OrderResponse{
		Symbol:           symbol,
		ClientOrderID:    uuid(10),
		ExecutedQuantity: executedQtyMock,
		Price:            fmt.Sprintf("%.15f", price),
	}
	// end DEBUG end mocking

	err = json.Unmarshal(response, orderResponse)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("orderResponse", orderResponse)
	price, errPrice := strconv.ParseFloat(orderResponse.Price, 64)
	executedQty, errQty := strconv.ParseFloat(orderResponse.ExecutedQuantity, 64)
	if errPrice != nil && errQty != nil {
		fmt.Println("errPrice", errPrice)
		fmt.Println("errQty", errPrice)
	}

	return &common.OrderResponse{
		Price:            price,
		QuantityAsked:    quantity,
		QuantityExecuted: executedQty,
		Platform:         a.Name(),
		OrderID:          orderResponse.ClientOrderID,
	}, nil
}

// DEBUG variable and functions

var HasFilled bool = false

func uuid(len int) string {
	buff := make([]byte, len)
	rand.Read(buff)
	str := base64.StdEncoding.EncodeToString(buff)
	// Base 64 can be longer than len
	return str[:len]
}

// END DEBUG
