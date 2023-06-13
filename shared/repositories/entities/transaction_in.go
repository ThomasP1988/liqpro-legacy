package entities

// StatusTransactionIn state of the executed transaction
type StatusTransactionIn string

// list of possible status
const (
	SuccessTransactionIn StatusTransactionIn = "success"
	FailedTransactionIn  StatusTransactionIn = "failed"
)

// TransactionIn Record of transaction between client and us
type TransactionIn struct {
	ID               string              `bson:"_id,omitempty" json:"id"`
	UserID           string              `bson:"userId" json:"userId"`
	Action           int                 `bson:"action" json:"action"`
	Instrument       string              `bson:"instrument" json:"instrument"`
	QuantityAsked    float64             `bson:"quantityAsked" json:"quantityAsked"`
	QuantityExecuted float64             `bson:"quantityexecuted" json:"quantityexecuted"`
	PayedUpfront     bool                `bson:"payedUpfront" json:"payedUpfront"`
	PricePerUnit     float64             `bson:"pricePerUnit" json:"pricePerUnit"`
	TotalPrice       float64             `bson:"totalPrice" json:"totalPrice"`
	Date             int64               `bson:"date" json:"dateExecuted"`
	Status           StatusTransactionIn `bson:"status" json:"status"`
	FailedReason     string              `bson:"failedReason" json:"failedReason"`
}
