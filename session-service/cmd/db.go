package main

import (
	"log"

	"github.com/adi290491/productivity-planner/session-service/config"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(appConfig *config.AppConfig) {

	db, err := gorm.Open(pg.Open(appConfig.DSN), &gorm.Config{})

	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	log.Printf("Database connection successful")
	appConfig.DB = db
}
