package models

import (
	"time"

	"github.com/google/uuid"
)

type UserDailyTrend struct {
	Id             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserId         uuid.UUID
	Day            time.Time
	FocusMinutes   float64
	MeetingMinutes float64
	BreakMinutes   float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type DailyTrendDao struct {
	UserId   string
	NoOfDays time.Time
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

type WeeklyTrendDao struct {
	UserId   string
	NoOfDays string
}
