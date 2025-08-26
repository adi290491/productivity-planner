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

	repo.FetchWeeklyTrend()

	log.Println("-------Job Execution Completed-------")
}
