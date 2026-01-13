package main

import (
	"log"

	"apipress/config"
	"apipress/database"
	"apipress/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load Config
	config.LoadConfig()

	// Connect to Database
	database.Connect()

	// Setup Gin
	r := gin.Default()

	// CORS Config (Basic)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup Routes
	routes.SetupRoutes(r)

	// Run Server
	port := config.AppConfig.Port
	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
