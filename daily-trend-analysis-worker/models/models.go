package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DailyAggregate struct {
	UserId         uuid.UUID `gorm:"column:user_id"`
	Day            time.Time `gorm:"column:day"`
	FocusMinutes   float64   `gorm:"column:focus_minutes"`
	MeetingMinutes float64   `gorm:"column:meeting_minutes"`
	BreakMinutes   float64   `gorm:"column:break_minutes"`
}

type UserDailyTrend struct {
	Id             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserId         uuid.UUID
	Day            time.Time
	FocusMinutes   float64
	MeetingMinutes float64
	BreakMinutes   float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ViewedAt       *time.Time
}

type Session struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	UserId      uuid.UUID
	SessionType string
	StartTime   time.Time
	EndTime     *time.Time
}

type UserBatch struct {
	UserIDs []uuid.UUID `json:"user_ids"`
}

type UserInfo struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

type UserSuccessInfo struct {
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
}

type UserFailureInfo struct {
	UserID uuid.UUID
	Errors error
}

type ProcessingSummary struct {
	SuccessfulUsers []UserSuccessInfo
	FailedUsers     []UserFailureInfo
	// FailureUserIDs  []uuid.UUID
	// Errors          []error
}

func (ps *ProcessingSummary) GetStats() string {
	return fmt.Sprintf("Processing Summary:\nSuccessful Users: %d\nFailed Users: %d", len(ps.SuccessfulUsers), len(ps.FailedUsers))
}

type TrendAnalysisEvent struct {
	Event           string            `json:"event"`
	JobType         string            `json:"job_type"`
	Status          string            `json:"status"`
	Date            string            `json:"date"`
	SuccessfulUsers []UserSuccessInfo `json:"successful_users,omitempty"` // Sent on success
	FailedUserIDs   []uuid.UUID       `json:"failed_user_ids,omitempty"`  // Sent on failure
	ErrorSummary    string            `json:"error_summary,omitempty"`    // Sent on failure
	NotifyAdmin     bool              `json:"notify_admin"`
}
