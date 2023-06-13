package repositories

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionUsers          = "users"
	collectionAPIKey         = "api_key"
	collectionTransactionIn  = "transaction_in"
	collectionTransactionOut = "transaction_out"
	collectionHold           = "hold"
	collectionToken          = "token"
)

var client *mongo.Database

//GetClient singleton like
func GetClient() *mongo.Database {
	if client == nil {
		ConnectDB()
	}
	return client
}

// ConnectDB connect to mongoDB
func ConnectDB() {

	ctex, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	clt, err := mongo.Connect(ctex, options.Client().ApplyURI(
		"mongodb+srv://thomas:5HVgsVKwwQ2xowPb@cluster0.0qzfq.mongodb.net/liqpro?retryWrites=true&w=majority",
	))
	if err != nil {
		log.Fatal(err)
	}
	db := clt.Database("liqpro")
	client = db
}
