package entities

// APIKey client api keys for rest api
type APIKey struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	AccessKey string `bson:"accessKey" json:"accessKey"`
	UserID    string `bson:"userId" json:"userId"`
	SecretKey string `bson:"secretKey" json:"secretKey"`
	Created   int64  `bson:"created" json:"created"`
}
