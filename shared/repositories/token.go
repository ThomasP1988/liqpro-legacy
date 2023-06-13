package repositories

import (
	"context"
	"fmt"
	"liqpro/shared/libs/crypto"
	"liqpro/shared/repositories/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TokenRep save all transaction requested by user
var tokenRep *TokenRepository

// GetTokenRepository singleton like function
func GetTokenRepository() *TokenRepository {
	if tokenRep == nil {
		tokenRep = &TokenRepository{
			client: GetClient(),
		}
	}
	return tokenRep
}

// TokenRepository to store token of apiKeys
type TokenRepository struct {
	client *mongo.Database
}

// Create token into DB
func (tr *TokenRepository) Create(userID string) (*entities.Token, error) {

	tokenValue, errGenerateString := crypto.GenerateString(20)

	if errGenerateString != nil {
		return nil, errGenerateString
	}

	token := &entities.Token{
		Value:  *tokenValue,
		UserID: userID,
	}

	insertResult, err := tr.client.Collection(collectionToken).InsertOne(context.TODO(), *token)
	if err != nil {
		return nil, err
	}

	fmt.Println("Inserted token with ID:", insertResult.InsertedID)
	return token, nil
}

func (tr *TokenRepository) Use(token string) (*string, error) {
	tknEntity := &entities.Token{}
	err := tr.FindOne(token, tknEntity)
	if err != nil {
		return nil, err
	}

	userID := tknEntity.UserID

	err = tr.Delete(token)

	if err != nil {
		fmt.Println("IMPORTANT: error deleting token", err)
	}

	return &userID, nil
}

// find primitive
func (tr *TokenRepository) find(filter *primitive.D, apiKey *entities.Token) error {
	err := tr.client.Collection(collectionToken).FindOne(context.TODO(), filter).Decode(apiKey)
	if err != nil {
		return err
	}

	fmt.Println("Found token for user ID:", apiKey.UserID)
	return nil
}

// FindOne user entity from DB
func (tr *TokenRepository) FindOne(token string, apiKey *entities.Token) error {

	return tr.find(&bson.D{{Key: "_id", Value: token}}, apiKey)
}

// Delete API key
func (tr *TokenRepository) Delete(token string) error {

	deleteResult, err := tr.client.Collection(collectionToken).DeleteOne(context.TODO(), &bson.D{{Key: "_id", Value: token}})
	if err != nil {
		return err
	}

	fmt.Println("Delete token result:", deleteResult)
	return nil
}
