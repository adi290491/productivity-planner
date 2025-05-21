package trend

import (
	"errors"
	"log"
	models "productivity-planner/trend-service/model"
	"productivity-planner/trend-service/utils"
	"time"
)

type TrendService struct {
	repo models.Repository
}

func NewTrendService(repo models.Repository) *TrendService {
	return &TrendService{
		repo: repo,
	}
}

func (t *TrendService) FetchDailyTrend(userId string, days string) (*DailyTrendResponse, error) {

	var noOfDays int

	noOfDays, ok := utils.ValidateDays(days)

	if !ok {
		return nil, errors.New("invalid value for 'days'. A positive integer is required")
	}

	dailyTrendDao := &models.DailyTrendDao{
		UserId:       userId,
		LookbackDays: time.Now().AddDate(0, 0, -noOfDays),
	}

	userDailyTrend, err := t.repo.FetchDailyTrend(dailyTrendDao)

	if err != nil {
		return nil, err
	}

	dailyTrendResponse := MapModelToResponse(userDailyTrend, userId)

	log.Println("Session Summary:", dailyTrendResponse)

	return dailyTrendResponse, nil
}
