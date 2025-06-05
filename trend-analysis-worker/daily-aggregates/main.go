package main

import (
	"log"
	"os"
)

func main() {

	app, err := LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(os.Stdout)
	log.Println("Daily trend job started")

	app.InitDB()

	repo := &PostgresRepository{
		DB: app.DB,
	}

	repo.FetchDailyTrends()
}
