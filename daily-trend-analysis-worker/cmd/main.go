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

	err = app.InitDB()

	if err != nil {
		log.Fatal(err)
	}

	repo := &PostgresRepository{
		DB: app.DB,
	}

	repo.FetchDailyTrends()
}
