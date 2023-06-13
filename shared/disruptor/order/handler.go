package order

import (
	"fmt"
	"time"

	common "liqpro/exchanges/common"
	repositories "liqpro/shared/repositories"
	"liqpro/shared/repositories/entities"

	disruptor "github.com/smartystreets-prototypes/go-disruptor"
)

// Define constant for the disruptor
const (
	BufferSize   = 32
	BufferMask   = BufferSize - 1
	Reservations = 1
)

var sequence int64 = 0
var ringBuffer = [BufferSize]*ToTrigger{}

// OrderDisruptor disruptor for orders
var OrderDisruptor = disruptor.New(
	disruptor.WithCapacity(BufferSize),
	disruptor.WithConsumerGroup(Consumer{}))

// Consumer struct we have to implement to use the library
type Consumer struct{}

// TriggerOrder trigger sell or buy on a platform
func TriggerOrder(userID []byte, requestID string, symbol string, price float64, quantity float64, side int, platform string, responseCh *(chan *common.OrderResponse)) {
	sequence = OrderDisruptor.Reserve(Reservations)

	ringBuffer[sequence&BufferMask] = &ToTrigger{
		UserID:          userID,
		RequestID:       requestID,
		Symbol:          symbol,
		Price:           price,
		Quantity:        quantity,
		Platform:        platform,
		Side:            side,
		ResponseChannel: responseCh,
	}

	OrderDisruptor.Commit(sequence-Reservations+1, sequence)
}

// Consume function to execute when thread is called
func (src Consumer) Consume(lower, upper int64) {
	fmt.Println("ici")
	for ; lower <= upper; lower++ {
		orderToTrigger := *ringBuffer[lower&BufferMask]

		platform := repositories.PlatformFactory.Get(orderToTrigger.Platform)

		var response *common.OrderResponse
		var err error

		if orderToTrigger.Side == 0 {
			response, err = (*platform).Buy(orderToTrigger.Symbol, orderToTrigger.Price, orderToTrigger.Quantity)
		} else {
			response, err = (*platform).Sell(orderToTrigger.Symbol, orderToTrigger.Price, orderToTrigger.Quantity)
		}

		if err != nil {
			// TODO: should tell the channel that a error happen for it to trigger more order
		}

		*orderToTrigger.ResponseChannel <- response
		fmt.Println("response", response)

		// save transaction out

		transactionOut := &entities.TransactionOut{
			UserID:           string(orderToTrigger.UserID),
			Platform:         orderToTrigger.Platform,
			Instrument:       orderToTrigger.Symbol,
			TransactionInID:  orderToTrigger.RequestID,
			Price:            response.Price,
			QuantityAsked:    orderToTrigger.Quantity,
			QuantityExecuted: response.QuantityExecuted,
			PlatformOrderID:  response.OrderID,
			Date:             time.Now().UnixNano(),
			Status:           entities.SuccessTransactionOut,
		}

		err = repositories.GetTransactionOutRepository().Create(transactionOut)

		if err != nil {
			fmt.Println("err", err)
		}

	}
}

// ToTrigger we send this struct to the disruptor to send an order
type ToTrigger struct {
	UserID          []byte
	RequestID       string
	Symbol          string
	Price           float64
	Quantity        float64
	Platform        string
	Side            int
	ResponseChannel *(chan *common.OrderResponse)
}
