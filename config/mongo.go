package config

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	mongoOnce   sync.Once
)

// GetMongoClient returns a singleton MongoDB client
func GetMongoClient() *mongo.Client {
	mongoOnce.Do(func() {
		uri := os.Getenv("MONGO_URI")
		if uri == "" {
			uri = "mongodb://localhost:27017"
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
		// Ping to ensure connection
		if err := client.Ping(ctx, nil); err != nil {
			log.Fatalf("MongoDB ping error: %v", err)
		}
		mongoClient = client
		log.Println("Connected to MongoDB")
	})
	return mongoClient
}
