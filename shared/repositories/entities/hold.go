package entities

import "liqpro/config"

// Hold keep information on how much the user holds for each currency
type Hold struct {
	ID           string          `bson:"_id" json:"id"`
	UserID       string          `bson:"userId" json:"userId"`
	Currency     config.Currency `bson:"currency" json:"currency"`
	Total        float64         `bson:"total" json:"total"`
	LastModified int             `bson:"lastModified" json:"lastModified"`
}
