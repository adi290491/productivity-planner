package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"time"

	"github.com/adi290491/productivity-planner/daily-trend-analysis-worker/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"net/http"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func (p *PostgresRepository) FetchDailyTrends(app *Application) (*models.ProcessingSummary, error) {

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
		return nil, fmt.Errorf("aggregation error: %v", err)
	}

	var userIDs []uuid.UUID
	for _, row := range dailyAggregate {
		userIDs = append(userIDs, row.UserId)
	}

	userInfo, err := p.fetchUserEmailsInBatch(app, userIDs)

	if err != nil {
		return nil, fmt.Errorf("error while fetching user emails: %+v", err)
	}

	summary := &models.ProcessingSummary{}

	log.Printf("Daily Aggregate: %+v", dailyAggregate)
	for _, row := range dailyAggregate {
		dailyTrend := models.UserDailyTrend{
			UserId:         row.UserId,
			Day:            row.Day,
			FocusMinutes:   row.FocusMinutes,
			MeetingMinutes: row.MeetingMinutes,
			BreakMinutes:   row.BreakMinutes,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			ViewedAt:       nil,
		}

		// log.Printf("Daily Trend: %+v\n", dailyTrend)
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

	// log.Printf("Processing complete. Success: %d, Failed: %d", len(summary.SuccessfulUsers), len(summary.FailedUsers))
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
	log.Printf("Users: %+v", string(users))
	if err != nil {
		log.Printf("[DEBUG] error reading response body: %v", err)
		return nil, err
	}

	resp.Body.Close()

	var userInfo []models.UserInfo
	err = json.Unmarshal(users, &userInfo)
	log.Printf("User Info: %+v", userInfo)
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
