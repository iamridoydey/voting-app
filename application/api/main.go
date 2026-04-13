package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"voting-app/api/db"
	"voting-app/api/handlers"
)

func main() {
	db.Connect()

	if err := db.SeedLanguages(); err != nil {
		// Non-fatal: log and continue — app is still usable if seed fails
		log.Printf("[main] Seed warning: %v", err)
	}

	r := gin.Default()

	// CORS: read from ConfigMap via env var FRONTEND_URLS
	// ConfigMap value: "http://localhost:5173,http://127.0.0.1:5173"
	origins := parseOrigins(os.Getenv("FRONTEND_URLS"))
	log.Printf("[main] CORS allowed origins: %v", origins)

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

	// Graceful shutdown: wait for SIGTERM (what Kubernetes sends on pod termination)
	// before closing the DB connection — in-flight requests finish cleanly
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Printf("[main] Server starting on port %s", port)
		if err := r.Run(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[main] Server error: %v", err)
		}
	}()

	<-quit
	log.Println("[main] Shutting down...")
	db.Disconnect()
	log.Println("[main] Shutdown complete")
}

func parseOrigins(raw string) []string {
	if raw == "" {
		return []string{"http://localhost:5173"}
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}