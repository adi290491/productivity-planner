package main

import (
	"productivity-planner/trend-analysis-worker/daily-aggregates/dailytrend"
	"productivity-planner/trend-analysis-worker/daily-aggregates/db"
)

func main() {
	db.InitDB()
	dailytrend.FetchDailyTrends()
}
