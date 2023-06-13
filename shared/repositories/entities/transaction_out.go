package entities

// StatusTransactionOut state of the executed transaction
type StatusTransactionOut string

// list of possible status
const (
	SuccessTransactionOut StatusTransactionOut = "success"
	FailedTransactionOut  StatusTransactionOut = "failed"
)

// TransactionOut record of a transaction between us and the exchanges
type TransactionOut struct {
	ID               string               `bson:"_id,omitempty"`
	UserID           string               `bson:"userId"`
	Platform         string               `bson:"platform"`
	Instrument       string               `bson:"instrument"`
	Price            float64              `bson:"price"`
	QuantityAsked    float64              `bson:"quantityAsked"`
	QuantityExecuted float64              `bson:"quantityExecuted"`
	TransactionInID  string               `bson:"transactionInId"`
	PlatformOrderID  string               `bson:"platformOrderID"`
	Date             int64                `bson:"date"`
	Status           StatusTransactionOut `bson:"status"`
}
