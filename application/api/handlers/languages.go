package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"voting-app/api/db"
	"voting-app/api/models"
)

func GetLanguages(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Sort by votes descending so the frontend always gets a ranked list
	opts := options.Find().SetSort(bson.D{{Key: "votes", Value: -1}})
	cursor, err := db.Collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch languages"})
		return
	}
	defer cursor.Close(ctx)

	var langs []models.Language
	if err := cursor.All(ctx, &langs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode languages"})
		return
	}

	c.JSON(http.StatusOK, langs)
}

func UpdateLanguage(c *gin.Context) {
	var body struct {
		Name  string `json:"name"`
		Delta int    `json:"delta"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Delta must be +1 or -1 — reject anything else to prevent vote manipulation
	if body.Delta != 1 && body.Delta != -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "delta must be 1 or -1"})
		return
	}
	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"name": body.Name}
	update := bson.M{"$inc": bson.M{"votes": body.Delta}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated models.Language
	if err := db.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "language not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}