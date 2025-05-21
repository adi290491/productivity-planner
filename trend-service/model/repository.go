package models

import "gorm.io/gorm"

type Repository interface {
	FetchDailyTrend(dailyTrendDao *UserDailyTrend) ([]UserDailyTrend, error)
	FetchWeeklyTrend(dailyTrendDao *UserWeeklyTrend) ([]UserWeeklyTrend, error)
}

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgresRepository(conn *gorm.DB) Repository {
	return &PostgresRepository{
		DB: conn,
	}
}
