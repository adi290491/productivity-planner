package main

import (
	"context"
	"log"
	"productivity-planner/trend-analysis-worker/daily-aggregates/models"

	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func (p *PostgresRepository) FetchDailyTrends() {

	db := p.DB

	log.Println("Fetching daily trends...")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var dailyAggregate []models.DailyAggregate
	err := db.WithContext(ctx).
		Model(&models.Session{}).
		Select(`
		user_id,
		date_trunc('day', start_time) AS day,
		SUM(CASE WHEN session_type = 'focus' THEN EXTRACT(EPOCH FROM (end_time - start_time)) / 60 ELSE 0 END) as focus_minutes,
		SUM(CASE WHEN session_type = 'meeting' THEN EXTRACT(EPOCH FROM (end_time - start_time)) / 60 ELSE 0 END) as meeting_minutes,
		SUM(CASE WHEN session_type = 'break' THEN EXTRACT(EPOCH FROM (end_time - start_time)) / 60 ELSE 0 END) as break_minutes
	`).
		Where("end_time IS NOT NULL AND date_trunc('day', start_time) = date_trunc('day', now())").
		Group("user_id, day").Find(&dailyAggregate).Error

	if err != nil {
		log.Fatalf("aggregation error: %v", err)
	}

	for _, row := range dailyAggregate {
		dailyTrend := models.UserDailyTrend{
			UserId:         row.UserId,
			Day:            row.Day,
			FocusMinutes:   row.FocusMinutes,
			MeetingMinutes: row.MeetingMinutes,
			BreakMinutes:   row.BreakMinutes,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		result := db.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "user_id"}, {Name: "day"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"focus_minutes":   gorm.Expr("EXCLUDED.focus_minutes"),
					"meeting_minutes": gorm.Expr("EXCLUDED.meeting_minutes"),
					"break_minutes":   gorm.Expr("EXCLUDED.break_minutes"),
					"updated_at":      time.Now(),
				}),
			}).Create(&dailyTrend)

		if result.Error != nil {
			log.Printf("failed to upsert for user: %v: %v", row.UserId, result.Error)
		}

		log.Println("Rows inserted:", result.RowsAffected)
	}
}
