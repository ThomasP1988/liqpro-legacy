package repositories

import (
	"context"
	"fmt"
	"liqpro/shared/repositories/entities"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var apiKeyRep *APIKeyRepository

// GetAPIKeyRepository singleton like function
func GetAPIKeyRepository() *APIKeyRepository {
	if apiKeyRep == nil {
		apiKeyRep = &APIKeyRepository{
			client: GetClient(),
		}
	}
	return apiKeyRep
}

// APIKeyRepository to store api key of apiKeys
type APIKeyRepository struct {
	client *mongo.Database
}

// Create api key into DB
func (ar *APIKeyRepository) Create(apiKey *entities.APIKey) error {
	insertResult, err := ar.client.Collection(collectionAPIKey).InsertOne(context.TODO(), *apiKey)
	if err != nil {
		return err
	}

	fmt.Println("Inserted apiKey with ID:", insertResult.InsertedID)
	return nil
}

// find primitive
func (ar *APIKeyRepository) find(filter *primitive.D, apiKey *entities.APIKey) error {
	err := ar.client.Collection(collectionAPIKey).FindOne(context.TODO(), filter).Decode(apiKey)
	if err != nil {
		return err
	}

	fmt.Println("Found apiKey for user ID:", apiKey.UserID)
	return nil
}

// FindOne user entity from DB
func (ar *APIKeyRepository) FindOne(accessKey string, apiKey *entities.APIKey) error {

	return ar.find(&bson.D{{Key: "accessKey", Value: accessKey}}, apiKey)
}

// ListByUser list api keys by user
func (ar *APIKeyRepository) ListByUser(userID string) (*[]entities.APIKey, error) {
	apiKeys := []entities.APIKey{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := ar.client.Collection(collectionAPIKey).Find(ctx, bson.D{{Key: "userId", Value: userID}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {

		apiKey := &entities.APIKey{}
		err := cur.Decode(apiKey)
		if err != nil {
			return nil, err
		}
		fmt.Println("apiKey", apiKey)

		// err = bson.Unmarshal(result, apiKey)
		apiKeys = append(apiKeys, *apiKey)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return &apiKeys, nil
}

// Delete API key
func (ar *APIKeyRepository) Delete(accessKey string) error {

	deleteResult, err := ar.client.Collection(collectionAPIKey).DeleteOne(context.TODO(), &bson.D{{Key: "accessKey", Value: accessKey}})
	if err != nil {
		return err
	}

	fmt.Println("Delete apiKey result:", deleteResult)
	return nil
}
