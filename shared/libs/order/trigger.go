package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"liqpro/config"
	commonExchanges "liqpro/exchanges/common"

	order "liqpro/shared/disruptor/order"
	repositories "liqpro/shared/repositories"
	"liqpro/shared/repositories/cache"
	entities "liqpro/shared/repositories/entities"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// Handle clients pass orderRequest here
func Handle(orderRequest *RequestArgs, userData *cache.UserDataCache) ([]byte, error) {

	fmt.Println("orderRequest", fmt.Sprintf("%.15f", orderRequest.Quantity))
	fmt.Println("instrument", orderRequest.Instrument)

	var requestID string = uuid.New().String()
	var prices []redis.Z
	var errPrices error
	if orderRequest.Side == "0" {
		prices, errPrices = repositories.GetBestPricesBuy(orderRequest.Instrument)

	} else if orderRequest.Side == "1" {
		prices, errPrices = repositories.GetBestPricesSell(orderRequest.Instrument)

	} else {
		return nil, errors.New("Wrong side")
	}

	if errPrices != nil {
		fmt.Print("errPrices", errPrices)
		action, _ := strconv.ParseInt(orderRequest.Side, 10, 8)

		transactionIn := &entities.TransactionIn{
			ID:           requestID,
			Action:       int(action),
			UserID:       string(userData.UserID),
			Instrument:   orderRequest.Instrument,
			Status:       entities.FailedTransactionIn,
			Date:         time.Now().UnixNano(),
			FailedReason: errPrices.Error(),
		}
		go repositories.GetTransactionInRepository().Create(transactionIn)
		return nil, errors.New("Error processing")
	}

	sideInt, _ := strconv.Atoi(orderRequest.Side)
	responseCh := make(chan *commonExchanges.OrderResponse, 1)

	CalculateAndTrigger(userData.UserID, requestID, sideInt, orderRequest, orderRequest.Quantity, &prices, &responseCh)

	fmt.Println("prices length", len(prices))

	var totalAcquired float64
	var pricePerUnit float64
	var numberOfAnswer float64
S:
	for {
		select {
		case orderReponse := <-responseCh:
			fmt.Println("received message", orderReponse)
			quantitySlipped := orderReponse.QuantityAsked - orderReponse.QuantityExecuted
			if quantitySlipped != 0 {
				if orderRequest.Side == "0" {
					prices, errPrices = repositories.GetBestPricesBuy(orderRequest.Instrument)

				} else if orderRequest.Side == "1" {
					prices, errPrices = repositories.GetBestPricesSell(orderRequest.Instrument)
				}
				CalculateAndTrigger(userData.UserID, requestID, sideInt, orderRequest, quantitySlipped, &prices, &responseCh)
			}
			numberOfAnswer++

			totalAcquired += orderReponse.QuantityExecuted
			pricePerUnit += orderReponse.Price

			if totalAcquired == orderRequest.Quantity {
				fmt.Println("finished")
				break S
			}

		default:
			// fmt.Println("no message received")
		}
	}
	fmt.Println("finished")
	pricePerUnit = pricePerUnit / numberOfAnswer
	// TODO: calculate our profit and add it to the price

	// response
	orderResponse := &RequestResponse{
		Event:        "orderResult",
		Price:        fmt.Sprintf("%.2f", pricePerUnit*totalAcquired),
		Instrument:   orderRequest.Instrument,
		Quantity:     totalAcquired,
		PricePerUnit: fmt.Sprintf("%.2f", pricePerUnit),
		Side:         orderRequest.Side,
		ClientID:     orderRequest.ClientID,
	}

	bson, _ := json.Marshal(orderResponse)

	SaveTransaction(orderRequest, userData, totalAcquired, pricePerUnit)

	return bson, nil

}

// CalculateAndTrigger calculate which platforms we need to send orders to
func CalculateAndTrigger(userID []byte, requestID string, sideInt int, orderRequest *RequestArgs, quantityNeeded float64, prices *[]redis.Z, responseCh *(chan *commonExchanges.OrderResponse)) {
	dataLen := len(*prices)
	var quantityProcessed float64 = 0
	var aggregateByPlatform map[string]float64 = map[string]float64{}
	var maxPriceByPlatform map[string]float64 = map[string]float64{}

	for i := 0; i < dataLen; i++ {
		dataArr := strings.Split((*prices)[i].Member.(string), ":")
		volume, err := strconv.ParseFloat(dataArr[0], 64)

		if err != nil {
			continue
		}
		// price prices[i].Score
		if volume+quantityProcessed >= quantityNeeded {
			quantityToProcess := quantityNeeded - quantityProcessed
			// process with quantityToProcess

			aggregateByPlatform[dataArr[2]] += quantityToProcess

			maxPriceByPlatform[dataArr[2]] = (*prices)[i].Score
			quantityProcessed += quantityToProcess
			break
		} else {
			// process with volume
			aggregateByPlatform[dataArr[2]] += volume
			maxPriceByPlatform[dataArr[2]] = (*prices)[i].Score
			quantityProcessed += volume
		}

	}

	for k, v := range aggregateByPlatform {
		order.TriggerOrder(userID, requestID, orderRequest.Instrument, maxPriceByPlatform[k], v, sideInt, k, responseCh)
	}
}

// SaveTransaction save the transaction done with the client
func SaveTransaction(orderRequest *RequestArgs, userData *cache.UserDataCache, totalAcquired float64, pricePerUnit float64) {

	action, _ := strconv.ParseInt(orderRequest.Side, 10, 8)
	totalPrice := pricePerUnit * totalAcquired
	transactionIn := &entities.TransactionIn{
		Action:           int(action),
		UserID:           string(userData.UserID),
		Instrument:       orderRequest.Instrument,
		Status:           entities.SuccessTransactionIn,
		Date:             time.Now().UnixNano(),
		QuantityAsked:    orderRequest.Quantity,
		PayedUpfront:     false,
		QuantityExecuted: totalAcquired,
		PricePerUnit:     pricePerUnit,
		TotalPrice:       totalPrice,
	}

	repositories.GetTransactionInRepository().Create(transactionIn)

	pairInfo := config.PairsInfo[orderRequest.Instrument]
	var incrementCurrency config.Currency
	var decrementCurrency config.Currency
	var incrementPrice float64
	var decrementPrice float64

	if action == 0 { // buy
		incrementCurrency = pairInfo.Base
		decrementCurrency = pairInfo.Quote
		incrementPrice = totalAcquired
		decrementPrice = totalPrice
	} else if action == 1 { // sell
		decrementCurrency = pairInfo.Base
		incrementCurrency = pairInfo.Quote
		incrementPrice = totalPrice
		decrementPrice = totalAcquired
	}

	repositories.GetPortfolioRepository().IncrementHold(string(userData.UserID), incrementCurrency, incrementPrice)
	repositories.GetPortfolioRepository().DecrementHold(string(userData.UserID), decrementCurrency, decrementPrice)

}
