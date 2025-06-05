package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadConfig() (*Application, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DB_DSN environment variable is required")
	}
	return &Application{
		DSN: dsn,
	}, nil
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
