// +build debug

package bitstamp

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	common "liqpro/exchanges/common"
	"net/url"
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
	params.Set("price", fmt.Sprint(0)) // prevent buying anything, bitstamp doesnt provide test environment
	params.Set("amount", fmt.Sprintf("%.15f", quantity))
	params.Set("ioc_order", "True")

	response, err := RequestSigned("POST", "/"+side+"/"+symbol+"/", params)
	fmt.Println("response", string(response))
	if err != nil {
		return nil, err
	}

	// DEBUG Mocking for debug purposes
	var executedQtyMock float64
	// we want to try when its fully filled or partially

	if HasFilled {
		executedQtyMock = quantity
	} else {
		executedQtyMock = quantity / 2
	}
	HasFilled = !HasFilled

	// end DEBUG end mocking

	result := &LimitOrderResult{
		Price:    price,
		Id:       time.Now().UnixNano() * 4,
		DateTime: time.Now().Local().String(),
		Amount:   executedQtyMock,
	}
	// err = json.Unmarshal(response, result)

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
