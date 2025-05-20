package main

import (
	"productivity-planner/trend-analysis-worker/weekly-aggregates/db"
	"productivity-planner/trend-analysis-worker/weekly-aggregates/weeklytrend"
)

func main() {
	db.InitDB()
	weeklytrend.FetchWeeklyTrend()
}
