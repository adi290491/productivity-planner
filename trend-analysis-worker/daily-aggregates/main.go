package main

import (
	"log"
	"os"
	"productivity-planner/trend-analysis-worker/daily-aggregates/dailytrend"
	"productivity-planner/trend-analysis-worker/daily-aggregates/db"
)

func main() {

	log.SetOutput(os.Stdout)
	log.Println("Daily trend job started")

	db.InitDB()
	dailytrend.FetchDailyTrends()
}
