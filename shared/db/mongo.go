package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// ConnectMongo initializes the MongoDB connection
func ConnectMongo(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	// Verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Println("Connected to MongoDB")
	Client = client
	return nil
}

// GetCollection retrieves a MongoDB collection
func GetCollection(database, collection string) *mongo.Collection {
	return Client.Database(database).Collection(collection)
}

func CreateUniqueIndex() error {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME not provided!")
	}

	usersCollection := Client.Database(dbName).Collection("users")
	emailIndex := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // Index on "email" field
		Options: options.Index().SetUnique(true),
	}

	waIndex := mongo.IndexModel{
		Keys:    bson.M{"phone_number": 1},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := usersCollection.Indexes().CreateOne(ctx, emailIndex)
	if err != nil {
		log.Printf("Failed to create index on email field: %v", err)
		return err
	}

	_, err = usersCollection.Indexes().CreateOne(ctx, waIndex)
	if err != nil {
		log.Printf("Failed to create index on phone_number field: %v", err)
		return err
	}

	log.Println("Unique index created on email field")
	log.Println("Unique index created on phone_number field")
	return nil
}
