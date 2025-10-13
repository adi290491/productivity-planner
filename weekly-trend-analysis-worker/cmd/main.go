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
	log.Println("Weekly trend job started")

	err = app.InitDB()

	if err != nil {
		log.Fatal(err)
	}

	repo := &PostgresRepository{
		DB: app.DB,
	}

	processingSummary, err := repo.FetchWeeklyTrend(app)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Daily Trend processed successfully:\n %v", processingSummary.GetStats())
	log.Println("-------Job Execution Completed-------")
}

