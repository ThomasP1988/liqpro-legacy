package coinbase

import (
	"encoding/json"
	"fmt"
	common "liqpro/exchanges/common"
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
	return "coinbase"
}

// Buy buy on coinbase
func (a PlatformAdapter) Buy(symbol string, price float64, quantity float64) (*common.OrderResponse, error) {
	fmt.Println("RELEASE BUY")
	return a.order("buy", symbol, price, quantity)
}

// Sell sell on coinbase
func (a PlatformAdapter) Sell(symbol string, price float64, quantity float64) (*common.OrderResponse, error) {
	fmt.Println("RELEASE SELL")
	return a.order("sell", symbol, price, quantity)
}

func (a PlatformAdapter) order(side string, symbol string, price float64, quantity float64) (*common.OrderResponse, error) {

	payloadBytes, err := json.Marshal(Order{
		Price:       strconv.FormatFloat(price, 'f', -1, 64),
		Size:        strconv.FormatFloat(quantity, 'f', -1, 64),
		Side:        side,
		ProductID:   pairsOursToThem[symbol],
		TimeInForce: "IOC",
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("price", strconv.FormatFloat(price, 'f', -1, 64))
	response, errReponse := RequestSigned("POST", "/orders", &payloadBytes)

	if errReponse != nil {
		return nil, errReponse
	}
	fmt.Println("response", string(response))
	orderResponse := &OrderResponse{}

	err = json.Unmarshal(response, orderResponse)

	if err != nil {
		return nil, err
	}

	price, errPrice := strconv.ParseFloat(orderResponse.ExecutedValue, 64)
	executedQty, errQty := strconv.ParseFloat(orderResponse.FilledSize, 64)
	if errPrice != nil && errQty != nil {
		fmt.Println("errPrice", errPrice)
		fmt.Println("errQty", errQty)
	}

	return &common.OrderResponse{
		Price:            price,
		QuantityAsked:    quantity,
		QuantityExecuted: executedQty,
		Platform:         a.Name(),
		OrderID:          orderResponse.ID,
	}, nil

}
