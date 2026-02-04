package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"voting-app/api/db"
	"voting-app/api/handlers"
)

func main() {
    // Connect DB + seed
    db.Connect()
    db.SeedLanguages()

    r := gin.Default()

    // Read env var FRONTEND_URLS (comma-separated list)
    frontendURLs := os.Getenv("FRONTEND_URLS")
    var origins []string
    if frontendURLs != "" {
        origins = strings.Split(frontendURLs, ",")
    } else {
        // fallback for local dev
        origins = []string{"http://localhost:5173"}
    }

    // Trim spaces just in case
    for i := range origins {
        origins[i] = strings.TrimSpace(origins[i])
    }

    r.Use(cors.New(cors.Config{
        AllowOrigins:     origins,
        AllowMethods:     []string{"GET", "POST"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
        AllowCredentials: true,
    }))

    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Welcome to voting-app API"})
    })

    r.GET("/languages", handlers.GetLanguages)
    r.POST("/languages", handlers.UpdateLanguage)

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }
    r.Run(":" + port)
}
