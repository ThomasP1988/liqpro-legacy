package repositories

import (
	"context"
	"fmt"
	"liqpro/config"
	"liqpro/shared/repositories/entities"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var portfolio *PortfolioRepository

// PortfolioRepository to store hold of currency
type PortfolioRepository struct {
	client *mongo.Database
}

// GetPortfolioRepository singleton like function
func GetPortfolioRepository() *PortfolioRepository {
	if portfolio == nil {
		portfolio = &PortfolioRepository{
			client: GetClient(),
		}
	}
	return portfolio
}

// Create hold into DB
func (pr *PortfolioRepository) Create(hold *entities.Hold) error {
	insertResult, err := pr.client.Collection(collectionHold).InsertOne(context.TODO(), *hold)
	if err != nil {
		return err
	}

	fmt.Println("Inserted apiKey with ID:", insertResult.InsertedID)
	return nil
}

// find primitive
func (pr *PortfolioRepository) find(filter *primitive.D, hold *entities.Hold) error {
	err := pr.client.Collection(collectionHold).FindOne(context.TODO(), filter).Decode(hold)
	if err != nil {
		return err
	}

	fmt.Println("Found apiKey for user ID:", hold.UserID)
	return nil
}

// FindOne user entity from DB
func (pr *PortfolioRepository) FindOne(id string, hold *entities.Hold) error {
	return pr.find(&bson.D{{Key: "_id", Value: id}}, hold)
}

// FindOne user entity from DB
func (pr *PortfolioRepository) Exists(id string) (bool, error) {
	hold := &entities.Hold{}

	err := pr.FindOne(id, hold)

	if err != nil {
		return false, err
	}

	if len(hold.ID) > 0 {
		return true, nil
	}

	return false, nil
}

// ListByUser list api keys by user
func (pr *PortfolioRepository) ListByUser(userID string) (*[]entities.Hold, error) {
	portfolio := []entities.Hold{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := pr.client.Collection(collectionHold).Find(ctx, bson.D{{Key: "userId", Value: userID}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {

		hold := &entities.Hold{}
		err := cur.Decode(hold)
		if err != nil {
			return nil, err
		}
		fmt.Println("hold", hold)

		// err = bson.Unmarshal(result, apiKey)
		portfolio = append(portfolio, *hold)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return &portfolio, nil
}

// Delete Hold
func (pr *PortfolioRepository) Delete(id string) error {
	deleteResult, err := pr.client.Collection(collectionHold).DeleteOne(context.TODO(), &bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}

	fmt.Println("Delete result:", deleteResult)
	return nil
}

func (pr *PortfolioRepository) crementHold(userID string, currency config.Currency, value float64) error {
	id := userID + string(currency)

	result, err := pr.client.Collection(collectionHold).UpdateOne(context.TODO(),
		&bson.M{
			"_id": id,
		},
		&bson.D{
			{Key: "$set", Value: bson.D{{Key: "lastModified", Value: time.Now().Unix()}}},
			{Key: "$inc", Value: bson.D{{Key: "total", Value: value}}},
			{Key: "$setOnInsert", Value: bson.D{
				{Key: "currency", Value: currency},
				{Key: "userId", Value: userID},
				// {Key: "total", Value: 0},
			}},
		},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		fmt.Println("err", err)
	}

	fmt.Println("result", result)
	return nil
}

// IncrementHold Hold
func (pr *PortfolioRepository) IncrementHold(userID string, currency config.Currency, value float64) error {
	return pr.crementHold(userID, currency, value)
}

// DecrementHold Hold
func (pr *PortfolioRepository) DecrementHold(userID string, currency config.Currency, value float64) error {
	return pr.crementHold(userID, currency, -value)
}
