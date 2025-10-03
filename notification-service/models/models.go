package models

import (
	"time"

	"github.com/google/uuid"
)

type TrendAnalysisEvent struct {
	Event           string           `json:"event"`
	JobType         string           `json:"jobType"`
	Status          string           `json:"status"`
	Date            string           `json:"date"`
	SuccessfulUsers []SuccessfulUser `json:"successfulUsers,omitempty"`
	FailedUserIDs   []string         `json:"failedUserIds,omitempty"`
	ErrorSummary    string           `json:"errorSummary,omitempty"`
	NotifyAdmin     bool             `json:"notifyAdmin"`
}

type SuccessfulUser struct {
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
}

type UserFailureInfo struct {
	UserID uuid.UUID
	Errors error
}

type UserNotification struct {
	ID                  uuid.UUID `gorm:"primaryKey"`
	UserID              uuid.UUID
	HasNewDailyTrend    bool
	LastDailyTrendDate  *time.Time
	LastDailyTrendID    *uuid.UUID
	HasNewWeeklyTrend   bool
	LastWeeklyTrendDate *time.Time
	LastWeeklyTrendID   *uuid.UUID
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (UserNotification) TableName() string {
	return "user_notifications"
}

type UserNotificationResponse struct {
	HasNewDailyTrend    bool       `json:"hasNewDailyTrend"`
	HasNewWeeklyTrend   bool       `json:"hasNewWeeklyTrend"`
	LastDailyTrendDate  *time.Time `json:"lastDailyTrendDate"`
	LastWeeklyTrendDate *time.Time `json:"lastWeeklyTrendDate"`
}
