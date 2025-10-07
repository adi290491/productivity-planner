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

	processingSummary, err := repo.FetchDailyTrends(app)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Daily Trend processed successfully:\n %v", processingSummary.GetStats())
	// log.Println("-------Daily Trend Processing Successful-------")

	// publisher := &Publisher{
	// 	ProcessingSummary: processingSummary,
	// }

	// publisher.Publish(context.Background(), app)

	log.Println("-------Job Execution Finished--------")

}
