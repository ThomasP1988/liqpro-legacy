// +build debug

package bittrex

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	common "liqpro/exchanges/common"
	crypto "liqpro/shared/libs/crypto"
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
		Limit:        0,
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
	id, _ := crypto.GenerateString(20)
	result := &OrderV3Response{
		ID:       *id,
		Limit:    price,
		Quantity: executedQtyMock,
	}

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
