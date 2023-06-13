package repositories

import (
	"context"
	"fmt"
	"liqpro/shared/repositories/entities"

	"go.mongodb.org/mongo-driver/mongo"
)

// TransactionOutRep instance of the repository
var transactionOutRep *TransactionOutRepository

// GetTransactionOutRepository singleton like function
func GetTransactionOutRepository() *TransactionOutRepository {
	if transactionOutRep == nil {
		transactionOutRep = &TransactionOutRepository{
			client: GetClient(),
		}
	}
	return transactionOutRep
}

// TransactionOutRepository to store api key of apiKeys
type TransactionOutRepository struct {
	client *mongo.Database
}

// Create api key into DB
func (ar *TransactionOutRepository) Create(transactionOut *entities.TransactionOut) error {
	insertResult, err := ar.client.Collection(collectionTransactionOut).InsertOne(context.TODO(), *transactionOut)
	if err != nil {
		return err
	}

	fmt.Println("Inserted transaction out with ID:", insertResult.InsertedID)
	return nil
}
