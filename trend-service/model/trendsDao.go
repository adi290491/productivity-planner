package models

import (
	"context"
	"fmt"
	"time"
)

func (p *PostgresRepository) FetchDailyTrend(dailyTrendDao *DailyTrendDao) ([]UserDailyTrend, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userDailyTrend []UserDailyTrend

	err := p.DB.WithContext(ctx).
		Where("user_id = ? AND day between ? AND CURRENT_DATE", dailyTrendDao.UserId, dailyTrendDao.NoOfDays).
		Find(&userDailyTrend).Error

	if len(userDailyTrend) == 0 {
		return nil, fmt.Errorf("no daily trends found for the last %v days", dailyTrendDao.NoOfDays)
	}

	if err != nil {
		return nil, err
	}
	return userDailyTrend, nil
}

func (p *PostgresRepository) FetchWeeklyTrend(dailyTrendDao *WeeklyTrendDao) ([]UserWeeklyTrend, error) {
	return nil, nil
}
