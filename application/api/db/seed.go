package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"voting-app/api/models"
)

// seedLanguages is the source of truth for initial data.
// To add/remove languages, edit this slice and redeploy.
var seedLanguages = []models.Language{
	{Name: "Rust", Votes: 0, Image: "https://www.rust-lang.org/logos/rust-logo-512x512.png"},
	{Name: "Go", Votes: 0, Image: "https://raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher.png"},
	{Name: "Python", Votes: 0, Image: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg"},
	{Name: "JavaScript", Votes: 0, Image: "https://upload.wikimedia.org/wikipedia/commons/6/6a/JavaScript-logo.png"},
	{Name: "Java", Votes: 0, Image: "https://upload.wikimedia.org/wikipedia/en/3/30/Java_programming_language_logo.svg"},
	{Name: "C++", Votes: 0, Image: "https://upload.wikimedia.org/wikipedia/commons/1/18/ISO_C%2B%2B_Logo.svg"},
	{Name: "C#", Votes: 0, Image: "https://upload.wikimedia.org/wikipedia/commons/4/4f/Csharp_Logo.png"},
	{Name: "PHP", Votes: 0, Image: "https://www.php.net/images/logos/php-logo.svg"},
	{Name: "TypeScript", Votes: 0, Image: "https://raw.githubusercontent.com/remojansen/logo.ts/master/ts.png"},
}

func SeedLanguages() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Unique index on "name" — this is what makes concurrent pod upserts safe.
	// If both replicas call SeedLanguages at the same time, MongoDB rejects the
	// duplicate on the second upsert attempt rather than creating duplicate docs.
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := Collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		// Non-fatal: index likely already exists from a previous pod startup
		log.Printf("[seed] Index warning (safe to ignore if already exists): %v", err)
	}

	inserted := 0
	for _, lang := range seedLanguages {
		filter := bson.M{"name": lang.Name}
		update := bson.M{
			// $setOnInsert: fields are ONLY written on a new document creation.
			// If the doc already exists (pod restart, second replica), nothing changes.
			// This means votes are never reset by a redeployment.
			"$setOnInsert": bson.M{
				"name":  lang.Name,
				"votes": lang.Votes,
				"image": lang.Image,
			},
		}
		result, err := Collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			log.Printf("[seed] Upsert failed for %s: %v", lang.Name, err)
			continue
		}
		if result.UpsertedCount > 0 {
			inserted++
		}
	}

	log.Printf("[seed] Done: %d new languages inserted", inserted)
	return nil
}