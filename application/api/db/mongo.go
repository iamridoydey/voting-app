package db

import (
    "context"
    "fmt"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Collection *mongo.Collection

func Connect() {
    // First check if a full URI is provided
    mongoURI := os.Getenv("MONGO_URI")

    if mongoURI == "" {
        // Otherwise build from individual parts
        user := os.Getenv("MONGO_USER")
        pass := os.Getenv("MONGO_PASSWORD")
        host := os.Getenv("MONGO_HOST") // EX: "localhost:27017" or "mongo:27017"
        dbName := os.Getenv("MONGO_DB") // EX: "votingdb"

        if host == "" {
            host = "localhost:27017"
        }
        if dbName == "" {
            dbName = "votingdb"
        }

        if user != "" && pass != "" {
            // Authenticated URI
            mongoURI = fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin", user, pass, host, dbName)
        } else {
            // No-auth fallback
            mongoURI = fmt.Sprintf("mongodb://%s/%s", host, dbName)
        }
    }

    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Connect to MongoDB
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        panic(err)
    }

    // Ping to verify connection
    if err := client.Ping(ctx, nil); err != nil {
        panic(err)
    }

    // Assign globals
    Client = client
    // Always use dbName from URI or env
    dbName := os.Getenv("MONGO_DB")
    if dbName == "" {
        dbName = "votingdb"
    }
    Collection = client.Database(dbName).Collection("languages")
}

// Optional: Graceful disconnect
func Disconnect() {
    if Client != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        if err := Client.Disconnect(ctx); err != nil {
            panic(err)
        }
    }
}
