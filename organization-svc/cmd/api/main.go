package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/api-suite/organization/src/router"
)

func main() {
	// Get JWT secret from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Get server port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Create Gin engine
	r := gin.Default()

	// Setup routes with JWT local authentication
	router.SetupRoutes(r, jwtSecret)

	// Start server
	log.Printf("Organization Service starting on port %s", port)
	log.Printf("Using JWT local validation (no remote calls)")

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
