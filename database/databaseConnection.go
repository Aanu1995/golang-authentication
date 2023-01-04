package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DatabaseInstance() *mongo.Client {
	// load environment variables
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	mongoDBURL := os.Getenv("MONGODB_URL")

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(mongoDBURL).SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()
	// connect to mongoDb
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connection success")

	return client
}

var Client *mongo.Client = DatabaseInstance()

func OpenCollection(collectionName string) *mongo.Collection {
	return Client.Database("Soarity").Collection(collectionName)
}