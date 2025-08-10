package main

import (
	"fmt"
	"log"

	"github.com/adi290491/productivity-planner/user-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(appConfig *config.AppConfig) error {
	log.Println("---------Initializing Database----------")
	// log.Printf("DSN: %s", appConfig.DSN)
	if appConfig == nil {
		return fmt.Errorf("InitDB: appConfig is nil")
	}

	if appConfig.DSN == "" {
		return fmt.Errorf("InitDB: DSN is empty")
	}
	db, err := gorm.Open(postgres.Open(appConfig.DSN), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	log.Printf("Database connection successful")

	appConfig.DB = db

	return nil
}
