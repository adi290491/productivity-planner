package main

import (
	"log"

	"github.com/adi290491/productivity-planner/notification-service/config"
	"github.com/adi290491/productivity-planner/notification-service/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(appConfig *config.AppConfig) {

	db, err := gorm.Open(postgres.Open(appConfig.DSN), &gorm.Config{})

	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	log.Printf("Database connection successful")

	// Auto-migrate the user notifications table
	if err := db.AutoMigrate(&models.UserNotification{}); err != nil {
		log.Fatalf("failed to migrate user notifications table: %v", err)
	}
	log.Printf("Database migration completed")

	appConfig.DB = db
}
