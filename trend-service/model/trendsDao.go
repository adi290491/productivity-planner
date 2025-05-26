package models

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func (p *PostgresRepository) FetchDailyTrend(dailyTrendDao *DailyTrendDao) ([]UserDailyTrend, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userDailyTrend []UserDailyTrend

	err := p.DB.WithContext(ctx).
		Where("user_id = ? AND day between ? AND CURRENT_DATE", dailyTrendDao.UserId, dailyTrendDao.LookbackDays).
		Find(&userDailyTrend).Error

	if len(userDailyTrend) == 0 {
		return nil, fmt.Errorf("no daily trends found for the last %v days", dailyTrendDao.LookbackDays)
	}

	if err != nil {
		return nil, err
	}
	return userDailyTrend, nil
}

func (p *PostgresRepository) FetchWeeklyTrend(weeklyTrendDao *WeeklyTrendDao) ([]UserWeeklyTrend, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userWeeklyTrend []UserWeeklyTrend

	err := p.DB.WithContext(ctx).
		Where("user_id = ? AND week_start between ? AND CURRENT_DATE", weeklyTrendDao.UserId, weeklyTrendDao.LookbackWeeks).
		Find(&userWeeklyTrend).Error

	if len(userWeeklyTrend) == 0 {
		return nil, fmt.Errorf("no weekly trends found for the last %v weeks", weeklyTrendDao.LookbackWeeks)
	}

	if err != nil {
		return nil, err
	}
	return userWeeklyTrend, nil
}
