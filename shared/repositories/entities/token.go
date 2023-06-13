package entities

type Token struct {
	Value  string `bson:"_id" json:"value"`
	UserID string `bson:"userId" json:"userId"`
}
