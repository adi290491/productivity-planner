package main

import (
	"log"
	"os"
	"productivity-planner/trend-analysis-worker/weekly-aggregates/db"
	"productivity-planner/trend-analysis-worker/weekly-aggregates/weeklytrend"
)

func main() {

	log.SetOutput(os.Stdout)
	log.Println("Daily trend job started")

	db.InitDB()
	weeklytrend.FetchWeeklyTrend()
}
