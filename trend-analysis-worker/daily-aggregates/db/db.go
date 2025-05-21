package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	var err error
	if db == nil {
		dsn := os.Getenv("DB_DSN")
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			log.Fatalf("init db error: %v", err)
		}
	}
}

func GetInstance() *gorm.DB {
	if db == nil {
		InitDB()
	}
	return db
}
