package main

import (
	"log"
	"productivity-planner/summary-service/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(appConfig *config.AppConfig) {

	db, err := gorm.Open(postgres.Open(appConfig.DSN), &gorm.Config{})

	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	log.Printf("Database connection successful")

	appConfig.DB = db
}
