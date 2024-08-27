package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateUniqueIndexes creates unique indexes for username and email fields.
func CreateUniqueIndexes(userCollection *mongo.Collection) error {
	ctx := context.Background()

	// Create a unique index on the username field
	usernameIndex := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	// Create a unique index on the email field
	emailIndex := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	// Create the indexes
	_, err := userCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{usernameIndex, emailIndex})
	if err != nil {
		return err
	}

	return nil
}
