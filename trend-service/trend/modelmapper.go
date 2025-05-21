package trend

import (
	"log"
	models "productivity-planner/trend-service/model"
)

func MapModelToResponse(dailyUserTrend []models.UserDailyTrend) *DailyTrendResponse {

	log.Println("Model object:", dailyUserTrend)
	return nil
}
