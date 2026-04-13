package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client
var Collection *mongo.Collection

func Connect() {
	// In production, MONGO_URI is the ONLY connection config you need.
	// For a StatefulSet ReplicaSet it looks like:
	// mongodb://user:pass@mongo-0.mongo:27017,mongo-1.mongo:27017,mongo-2.mongo:27017/votingdb?replicaSet=rs0&authSource=admin
	//
	// For local dev without auth:
	// mongodb://localhost:27017/votingdb
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		// Sensible local dev fallback — no URI building logic in production
		mongoURI = "mongodb://localhost:27017/votingdb"
		log.Println("[db] MONGO_URI not set, falling back to localhost")
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "votingdb"
	}

	// clientOptions: driver reads replicaSet, auth, TLS etc. directly from the URI
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		// For a ReplicaSet: prefer reading from secondary to reduce primary load.
		// Use readpref.Primary() if you need strong consistency on reads.
		SetReadPreference(readpref.SecondaryPreferred())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("[db] Failed to create mongo client: %v", err)
	}

	// Ping verifies at least one replicaset member is reachable
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("[db] Failed to ping MongoDB: %v", err)
	}

	Client = client
	Collection = client.Database(dbName).Collection("languages")
	log.Printf("[db] Connected to MongoDB | database: %s", dbName)
}

func Disconnect() {
	if Client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := Client.Disconnect(ctx); err != nil {
		log.Printf("[db] Error during disconnect: %v", err)
	}
}