package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"time"

	"github.com/adi290491/productivity-planner/weekly-trend-analysis-worker/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func (p *PostgresRepository) FetchWeeklyTrend(app *Application) (*models.ProcessingSummary, error) {

	log.Println("Fetching weekly trends...")

	db := p.DB

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var weeklyAggregate []models.WeeklyAggregate

	err := db.WithContext(ctx).
		Model(&models.Session{}).
		Select(`
			user_id,
			date_trunc('week', start_time) AS week_start,
			SUM(CASE WHEN session_type = 'focus' THEN EXTRACT(EPOCH FROM (end_time - start_time)) / 60 ELSE 0 END) as focus_minutes,
			SUM(CASE WHEN session_type = 'meeting' then EXTRACT(EPOCH FROM (end_time - start_time)) / 60 ELSE 0 END) as meeting_minutes,
			SUM(CASE WHEN session_type = 'break' then EXTRACT(EPOCH FROM (end_time - start_time)) / 60 ELSE 0 END) as break_minutes
	`).Where("end_time IS NOT NULL AND date_trunc('week', start_time) = date_trunc('week', now())").
		Group("user_id, week_start").Find(&weeklyAggregate).Error

	if err != nil {
		log.Fatalf("aggregation error: %v", err)
	}

	var userIDs []uuid.UUID
	for _, row := range weeklyAggregate {
		userIDs = append(userIDs, row.UserId)
	}

	userInfo, err := p.fetchUserEmailsInBatch(app, userIDs)

	if err != nil {
		return nil, fmt.Errorf("error while fetching user emails: %+v", err)
	}

	log.Printf("Weekly Aggregation: %+v", weeklyAggregate)

	summary := &models.ProcessingSummary{}

	for _, row := range weeklyAggregate {
		weeklyTrend := models.UserWeeklyTrend{
			UserId:         row.UserId,
			WeekStart:      row.WeekStart,
			FocusMinutes:   row.FocusMinutes,
			MeetingMinutes: row.MeetingMinutes,
			BreakMinutes:   row.BreakMinutes,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			ViewedAt:       nil,
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
			log.Printf("failed to upsert for user: %v: %v", row.UserId, result.Error)
			summary.FailedUsers = append(summary.FailedUsers, models.UserFailureInfo{
				UserID: row.UserId,
				Errors: result.Error,
			})
		} else {
			// Include both user ID and email for successful users
			userEmail := "" // fallback
			if userObj, exists := userInfo[row.UserId]; exists {
				userEmail = userObj.Email
			}

			summary.SuccessfulUsers = append(summary.SuccessfulUsers, models.UserSuccessInfo{
				UserID: row.UserId,
				Email:  userEmail,
			})
		}

		log.Println("Rows inserted:", result.RowsAffected)
	}
	return summary, nil
}

func (p *PostgresRepository) fetchUserEmailsInBatch(app *Application, userIDs []uuid.UUID) (map[uuid.UUID]models.UserInfo, error) {

	userBatch := models.UserBatch{
		UserIDs: userIDs,
	}

	buf, err := json.Marshal(userBatch)

	if err != nil {
		return nil, fmt.Errorf("failed to marshall userIds. %+v", err)
	}

	resp, err := http.Post(app.USER_BATCH_URL,
		"application/json",
		bytes.NewBuffer(buf),
	)

	if err != nil {
		log.Printf("[DEBUG] error from user service: %v", err)
		return nil, err
	}

	users, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Printf("[DEBUG] error reading response body: %v", err)
		return nil, err
	}

	resp.Body.Close()

	var userInfo []models.UserInfo
	err = json.Unmarshal(users, &userInfo)

	if err != nil {
		log.Printf("[DEBUG] error unmarshalling response body: %v", err)
		return nil, err
	}

	userMap := make(map[uuid.UUID]models.UserInfo)
	for _, user := range userInfo {
		userMap[user.Id] = user
	}

	return userMap, nil
}
