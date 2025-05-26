package trend

import (
	"log"
	models "productivity-planner/trend-service/model"
	"productivity-planner/trend-service/utils"
	"time"
)

type TrendService struct {
	Repo models.Repository
}

func NewTrendService(repo models.Repository) *TrendService {
	return &TrendService{
		Repo: repo,
	}
}

func (t *TrendService) FetchDailyTrend(userId string, days string) (*DailyTrendResponse, error) {

	var noOfDays int

	noOfDays, err := utils.ValidateDays(days)

	if err != nil {
		return nil, err
	}
	log.Println("No of days:", noOfDays)

	dailyTrendDao := &models.DailyTrendDao{
		UserId:       userId,
		LookbackDays: time.Now().AddDate(0, 0, -noOfDays),
	}

	userDailyTrend, err := t.Repo.FetchDailyTrend(dailyTrendDao)

	if err != nil {
		return nil, err
	}

	dailyTrendResponse := MapModelToResponse(userDailyTrend, userId)

	return dailyTrendResponse, nil
}

func (t *TrendService) FetchWeeklyTrend(userId string, weeks string) (*WeeklyTrendResponse, error) {

	var noOfWeeks int

	noOfWeeks, err := utils.ValidateDays(weeks)

	if err != nil {
		return nil, err
	}

	log.Println("No of weeks:", noOfWeeks)

	weeklyTrendDao := &models.WeeklyTrendDao{
		UserId:        userId,
		LookbackWeeks: time.Now().AddDate(0, 0, -noOfWeeks*7),
	}

	userWeeklyTrend, err := t.Repo.FetchWeeklyTrend(weeklyTrendDao)

	if err != nil {
		return nil, err
	}

	weeklyTrendResponse := MapWeeklyModelToResponse(userWeeklyTrend, userId)

	return weeklyTrendResponse, nil
}
