package models

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func (p *PostgresRepository) CountUnviewedDailyTrends(userID uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int64

	err := p.DB.WithContext(ctx).
		Model(&UserDailyTrend{}).Where("user_id = ? AND viewed_at IS NULL", userID).
		Count(&count).Error

	return int(count), err
}

func (p *PostgresRepository) CountUnviewedWeeklyTrends(userID uuid.UUID) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int64

	err := p.DB.WithContext(ctx).
		Model(&UserWeeklyTrend{}).Where("user_id = ? AND viewed_at IS NULL", userID).
		Count(&count).Error

	return int(count), err
}

func (p *PostgresRepository) MarkDailyTrendsAsViewed(userId string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return p.DB.WithContext(ctx).
        Model(&UserDailyTrend{}).
        Where("user_id = ? AND viewed_at IS NULL", userId).
        Update("viewed_at", time.Now()).Error
}

func (p *PostgresRepository) MarkWeeklyTrendsAsViewed(userId string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return p.DB.WithContext(ctx).
        Model(&UserWeeklyTrend{}).
        Where("user_id = ? AND viewed_at IS NULL", userId).
        Update("viewed_at", time.Now()).Error
}