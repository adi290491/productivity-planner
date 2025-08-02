package main

import (
	"log"

	"github.com/adi290491/productivity-planner/user-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(appConfig *config.AppConfig) {
	log.Printf("DSN: %s", appConfig.DSN)
	db, err := gorm.Open(postgres.Open(appConfig.DSN), &gorm.Config{})

	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	log.Printf("Database connection successful")

	appConfig.DB = db
}
