package main

import (
	"log"
	"os"

	"gorm.io/gorm"
)

type Application struct {
	DB  *gorm.DB
	DSN string
}

func main() {
	app, err := LoadConfig()

	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(os.Stdout)
	log.Println("Weekly trend job started")

	app.InitDB()
	repo := &PostgresRepository{
		DB: app.DB,
	}

	repo.FetchWeeklyTrend()
}
