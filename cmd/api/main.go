package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	// Initialize Gin with default middleware (logger and recovery)
	r := gin.Default()

	// Pass db to routes if needed (example: api.RegisterRoutes(apiGroup, db))
	// apiGroup := r.Group("/api")

	// Example: Register routes from other modules
	// invoice.RegisterRoutes(apiGroup, db)
	// user.RegisterRoutes(apiGroup, db)

	// Start server
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
