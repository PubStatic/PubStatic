package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var client mongo.Client

func setup(connectionString string) {
	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	tempClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	client = *tempClient

	// Ping the MongoDB server to verify that we're connected
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Connected to MongoDB!")
}

func WriteMongo(database string, collectionName string, document interface{}, connectionString string) error {
	setup(connectionString)
	// Disconnect from MongoDB when program exits
	defer func() {
		if discErr := client.Disconnect(context.Background()); discErr != nil {
			log.Fatal(discErr)
		}
	}()

	collection := client.Database(database).Collection(collectionName)

	// Insert a document
	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return err
	}

	logger.Info("Inserted document")

	return nil
}

func ReadMongo[T any](database string, collectionName string, filter interface{}, connectionString string) (T, error) {
	setup(connectionString)
	// Disconnect from MongoDB when program exits
	defer func() {
		if discErr := client.Disconnect(context.Background()); discErr != nil {
			log.Fatal(discErr)
		}
	}()

	collection := client.Database(database).Collection(collectionName)

	var result T

	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
