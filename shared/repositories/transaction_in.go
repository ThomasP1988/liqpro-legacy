package repositories

import (
	"context"
	"fmt"
	"liqpro/shared/repositories/entities"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransactionInRep save all transaction requested by user
var transactionInRep *TransactionInRepository

// GetTransactionInRepository singleton like function
func GetTransactionInRepository() *TransactionInRepository {
	if transactionInRep == nil {
		transactionInRep = &TransactionInRepository{
			client: GetClient(),
		}
	}
	return transactionInRep
}

// TransactionInRepository to store api key of apiKeys
type TransactionInRepository struct {
	client *mongo.Database
}

// Create api key into DB
func (tr *TransactionInRepository) Create(transactionIn *entities.TransactionIn) error {
	insertResult, err := tr.client.Collection(collectionTransactionIn).InsertOne(context.TODO(), *transactionIn)
	if err != nil {
		return err
	}

	fmt.Println("Inserted transaction in with ID:", insertResult.InsertedID)
	return nil
}

// ListByUser list api keys by user
func (tr *TransactionInRepository) ListByUser(userID string, skip *int64, limit *int64) (*[]entities.TransactionIn, error) {
	transactionsIn := []entities.TransactionIn{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := &options.FindOptions{
		Skip:  skip,
		Limit: limit,
	}

	cur, err := tr.client.Collection(collectionTransactionIn).Find(ctx, bson.D{{Key: "userId", Value: userID}}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {

		transactionIn := &entities.TransactionIn{}
		err := cur.Decode(transactionIn)
		if err != nil {
			return nil, err
		}
		fmt.Println("hold", transactionIn)

		// err = bson.Unmarshal(result, apiKey)
		transactionsIn = append(transactionsIn, *transactionIn)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return &transactionsIn, nil
}
