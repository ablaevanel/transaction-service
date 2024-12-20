package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
)

var Client *mongo.Client
var UsersCollection *mongo.Collection

func init() {
	if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://host.docker.internal:27017" 
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	Client = client

	UsersCollection = client.Database("auth_service").Collection("users")
	log.Println("Successfully connected to MongoDB!")
}

func Disconnect() {
	if err := Client.Disconnect(context.Background()); err != nil {
		log.Fatal("Failed to disconnect from MongoDB:", err)
	}
}
