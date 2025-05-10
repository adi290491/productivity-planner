package repository

import (
	"log"
	"productivity-planner/task-service/config"

	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(appConfig *config.AppConfig) {

	db, err := gorm.Open(pg.Open(appConfig.DSN), &gorm.Config{})

	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	log.Printf("Database connection successful. %v", db)
	appConfig.DB = db
}
