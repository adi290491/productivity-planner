package models

import "gorm.io/gorm"

type Repository interface {
	FetchDailyTrend(dailyTrendDao *DailyTrendDao) ([]UserDailyTrend, error)
	FetchWeeklyTrend(dailyTrendDao *WeeklyTrendDao) ([]UserWeeklyTrend, error)
}

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgresRepository(conn *gorm.DB) Repository {
	return &PostgresRepository{
		DB: conn,
	}
}
