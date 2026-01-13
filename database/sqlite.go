package database

import (
	"apipress/config"
	"apipress/models"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	// Ensure storage directory exists
	dir := filepath.Dir(config.AppConfig.DBPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal("Failed to create storage directory:", err)
	}

	DB, err = gorm.Open(sqlite.Open(config.AppConfig.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to SQLite database")

	// Auto Migration
    // Note: models package will be created in next steps, so we might need to update this file or ensure models exists first.
    // For now I will comment out the migration call until models are ready, or just include it and ensure I create models in this batch.
    // I will include it and ensure models are created in parallel.
	err = DB.AutoMigrate(&models.User{}, &models.Post{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
    log.Println("Database migration completed")
}
