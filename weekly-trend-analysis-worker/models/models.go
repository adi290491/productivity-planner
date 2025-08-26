package models

import (
	"time"

	"github.com/google/uuid"
)

type WeeklyTrendResult struct {
	UserId         uuid.UUID `gorm:"column:user_id"`
	WeekStart      time.Time `gorm:"column:week_start"`
	FocusMinutes   float64   `gorm:"column:focus_minutes"`
	MeetingMinutes float64   `gorm:"column:meeting_minutes"`
	BreakMinutes   float64   `gorm:"column:break_minutes"`
}

type UserWeeklyTrend struct {
	Id             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserId         uuid.UUID
	WeekStart      time.Time
	FocusMinutes   float64
	MeetingMinutes float64
	BreakMinutes   float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Session struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	UserId      uuid.UUID
	SessionType string
	StartTime   time.Time
	EndTime     *time.Time
}
