package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestDBRepo struct {
	DB *gorm.DB
}

func (p *TestDBRepo) FetchDailyTrend(dailyTrendDao *DailyTrendDao) ([]UserDailyTrend, error) {

	if dailyTrendDao == nil {
		return []UserDailyTrend{}, nil
	}
	return []UserDailyTrend{
		{
			UserId:         uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			Day:            time.Now().AddDate(0, 0, -1),
			FocusMinutes:   60,
			MeetingMinutes: 30,
			BreakMinutes:   15,
		},
	}, nil
}

func (p *TestDBRepo) FetchWeeklyTrend(weeklyTrendDao *WeeklyTrendDao) ([]UserWeeklyTrend, error) {

	return []UserWeeklyTrend{
		{
			UserId:         uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			WeekStart:      time.Now().AddDate(0, 0, -7),
			FocusMinutes:   240,
			MeetingMinutes: 90,
			BreakMinutes:   30,
		},
	}, nil
}
