package main

import (
	"context"
	"log"
	"productivity-planner/trend-analysis-worker/weekly-aggregates/models"

	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func (p *PostgresRepository) FetchWeeklyTrend() {

	log.Println("Fetching weekly trends...")

	db := p.DB

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var weeklyResults []models.WeeklyTrendResult

	err := db.WithContext(ctx).
		Model(&models.Session{}).
		Select(`
			user_id,
			date_trunc('week', start_time) AS week_start,
			SUM(CASE WHEN session_type = 'focus' THEN EXTRACT(EPOCH FROM (end_time - start_time)) / 60 ELSE 0 END) as focus_minutes,
			SUM(CASE WHEN session_type = 'meeting' then EXTRACT(EPOCH FROM (end_time - start_time)) / 60 ELSE 0 END) as meeting_minutes,
			SUM(CASE WHEN session_type = 'break' then EXTRACT(EPOCH FROM (end_time - start_time)) / 60 ELSE 0 END) as break_minutes
	`).Where("end_time IS NOT NULL AND date_trunc('week', start_time) = date_trunc('week', now())").
		Group("user_id, week_start").Find(&weeklyResults).Error

	if err != nil {
		log.Fatalf("aggregation error: %v", err)
	}

	for _, row := range weeklyResults {
		weeklyTrend := models.UserWeeklyTrend{
			UserId:         row.UserId,
			WeekStart:      row.WeekStart,
			FocusMinutes:   row.FocusMinutes,
			MeetingMinutes: row.MeetingMinutes,
			BreakMinutes:   row.BreakMinutes,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		result := db.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "user_id"}, {Name: "week_start"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"focus_minutes":   gorm.Expr("EXCLUDED.focus_minutes"),
					"meeting_minutes": gorm.Expr("EXCLUDED.meeting_minutes"),
					"break_minutes":   gorm.Expr("EXCLUDED.break_minutes"),
					"updated_at":      time.Now(),
				}),
			}).Create(&weeklyTrend)

		if result.Error != nil {
			log.Printf("failed to upsert for user %v: %v", row.UserId, result.Error)
		}

		log.Println("Rows inserted:", result.RowsAffected)
	}

}
