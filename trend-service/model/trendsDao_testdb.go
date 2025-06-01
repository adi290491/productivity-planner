package models

import (
	"github.com/google/uuid"
)

type TestDBRepo struct {
}

func (p *TestDBRepo) FetchDailyTrend(dailyTrendDao *DailyTrendDao) ([]UserDailyTrend, error) {

	return []UserDailyTrend{
		{
			UserId: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
	}, nil
}

func (p *TestDBRepo) FetchWeeklyTrend(weeklyTrendDao *WeeklyTrendDao) ([]UserWeeklyTrend, error) {

	return []UserWeeklyTrend{
		{
			UserId: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
	}, nil
}
