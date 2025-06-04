package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadConfig() *Application {
	return &Application{
		DSN: os.Getenv("DB_DSN"),
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
