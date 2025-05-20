package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	var err error
	if db == nil {
		dsn := "postgres://adityasawant:S10dulkar@localhost:5433/productivity_planner?sslmode=disable"
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
