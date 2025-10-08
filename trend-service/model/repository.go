package models

import "github.com/google/uuid"

type Repository interface {
	FetchDailyTrend(dailyTrendDao *DailyTrendDao) ([]UserDailyTrend, error)
	FetchWeeklyTrend(dailyTrendDao *WeeklyTrendDao) ([]UserWeeklyTrend, error)
	CountUnviewedDailyTrends(userID uuid.UUID) (int, error)
	CountUnviewedWeeklyTrends(userID uuid.UUID) (int, error)
	MarkDailyTrendsAsViewed(userId string) error
	MarkWeeklyTrendsAsViewed(userId string) error
}
