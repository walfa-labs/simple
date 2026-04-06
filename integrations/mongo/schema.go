package mongo

import (
	"context"

	"simple/constants"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MigrateSchema creates the required indexes for MongoDB
// MongoDB is schemaless, but indexes improve query performance
func MigrateSchema(collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	// Create index on created_at field for efficient sorting
	indexModel := mongo.IndexModel{
		Keys:    bson.D{bson.E{Key: "created_at", Value: -1}},
		Options: options.Index().SetName("idx_created_at"),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		// Index already exists is not an error
		if err.Error() != "Index: idx_created_at already exists" {
			return err
		}
	}

	return nil
}

// GetSchemaVersion returns the current schema version
func GetSchemaVersion() string {
	return "v1"
}