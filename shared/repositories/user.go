package repositories

import (
	"context"
	"fmt"
	entities "liqpro/shared/repositories/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userRep *UserRepository

// GetUserRepository singleton like function
func GetUserRepository() *UserRepository {
	if userRep == nil {
		userRep = &UserRepository{
			client: GetClient(),
		}
	}
	return userRep
}

// UserRepository repository to handle users
type UserRepository struct {
	client *mongo.Database
}

// Create user entity to DB
func (ur *UserRepository) Create(user *entities.User) error {
	insertResult, err := ur.client.Collection(collectionUsers).InsertOne(context.TODO(), *user)
	if err != nil {
		return err
	}

	fmt.Println("Inserted user with ID:", insertResult.InsertedID)
	return nil
}

// Replace user entity to DB
func (ur *UserRepository) Replace(id string, user *entities.User) error {

	opts := options.Replace().SetUpsert(false)

	replaceResult, err := ur.client.Collection(collectionUsers).ReplaceOne(context.TODO(), bson.D{{Key: "_id", Value: id}}, user, opts)
	if err != nil {
		return err
	}

	fmt.Println("Inserted user with ID:", replaceResult.UpsertedID)
	return nil
}

func (ur *UserRepository) find(filter *primitive.D, user *entities.User) error {

	err := ur.client.Collection(collectionUsers).FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		return err
	}

	fmt.Println("Found user with ID:", user.ID)
	return nil
}

// FindOneByMail user entity from DB
func (ur *UserRepository) FindOneByMail(email string, user *entities.User) error {
	return ur.find(&bson.D{{Key: "email", Value: email}}, user)
}

// FindOne user entity from DB
func (ur *UserRepository) FindOne(id string, user *entities.User) error {
	return ur.find(&bson.D{{Key: "_id", Value: id}}, user)
}

// FindOneByConfirmSelector user entity from DB by Confirm Selector
func (ur *UserRepository) FindOneByConfirmSelector(confirmselector string, user *entities.User) error {
	return ur.find(&bson.D{{Key: "confirmselector", Value: confirmselector}}, user)
}

// FindOneByRecoverSelector user entity from DB by Confirm Selector
func (ur *UserRepository) FindOneByRecoverSelector(recoverselector string, user *entities.User) error {
	return ur.find(&bson.D{{Key: "recoverselector", Value: recoverselector}}, user)
}
