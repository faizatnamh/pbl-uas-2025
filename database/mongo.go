package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// Init MongoDB Connection
func ConnectMongo() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is missing in .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect MongoDB:", err)
	}

	MongoClient = client
	log.Println("✓ MongoDB connected")
}

// Helper untuk ambil DB
func MongoDB() *mongo.Database {
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "uasbackend" // default name sesuai screenshot kamu
	}
	return MongoClient.Database(dbName)
}

// Helper untuk ambil Collection
func MongoCollection(name string) *mongo.Collection {
	return MongoDB().Collection(name)
}
