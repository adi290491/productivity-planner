package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadConfig() *Application {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}
	return &Application{
		DSN: dsn,
	}
}

func (a *Application) InitDB() {

	db, err := gorm.Open(postgres.Open(a.DSN), &gorm.Config{})

	if err != nil {
		log.Fatalf("init db error: %v", err)
	}

	sqldb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err = sqldb.Ping(); err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	a.DB = db
}
