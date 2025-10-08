package trend

import (
	"log"

	"time"

	models "github.com/adi290491/productivity-planner/trend-service/model"
	"github.com/adi290491/productivity-planner/trend-service/utils"
	"github.com/google/uuid"
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

func (t *TrendService) GetUnviewedTrendsCount(userId string) (*UnviewedTrendsCount, error) {

	userID, err := uuid.Parse(userId)

	if err != nil {
		return nil, err
	}

	dailyCount, err := t.Repo.CountUnviewedDailyTrends(userID)
	if err != nil {
		return nil, err
	}

	weeklyCount, err := t.Repo.CountUnviewedWeeklyTrends(userID)
	if err != nil {
		return nil, err
	}

	return &UnviewedTrendsCount{
		DailyCount:  dailyCount,
		WeeklyCount: weeklyCount,
	}, nil
}

func (t *TrendService) MarkTrendsAsViewed(userId string, trendType string) error {
	if trendType == "daily" {
		return t.Repo.MarkDailyTrendsAsViewed(userId)
	}
	return t.Repo.MarkWeeklyTrendsAsViewed(userId)
}
