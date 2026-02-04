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

    cursor, err := db.Collection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    var langs []models.Language
    if err := cursor.All(ctx, &langs); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, langs)
}

func UpdateLanguage(c *gin.Context) {
    var body struct {
        Name  string `json:"name"`
        Delta int    `json:"delta"`
    }
    if err := c.BindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"name": body.Name}
    update := bson.M{"$inc": bson.M{"votes": body.Delta}}

    result := db.Collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
    var updated models.Language
    if err := result.Decode(&updated); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "language not found"})
        return
    }
    c.JSON(http.StatusOK, updated)
}
