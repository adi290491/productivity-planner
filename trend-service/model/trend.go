package models

import (
	"fmt"
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
	UserId       string
	LookbackDays time.Time
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

func (u UserDailyTrend) String() string {
	return fmt.Sprintf(
		"UserDailyTrend{Id: %s, UserId: %s, Day: %s, FocusMinutes: %.2f, MeetingMinutes: %.2f, BreakMinutes: %.2f, CreatedAt: %s, UpdatedAt: %s}",
		u.Id, u.UserId, u.Day.Format(time.RFC3339), u.FocusMinutes, u.MeetingMinutes, u.BreakMinutes, u.CreatedAt.Format(time.RFC3339), u.UpdatedAt.Format(time.RFC3339),
	)
}

func (d DailyTrendDao) String() string {
	return fmt.Sprintf(
		"DailyTrendDao{UserId: %s, NoOfDays: %s}",
		d.UserId, d.LookbackDays.Format(time.RFC3339),
	)
}

func (u UserWeeklyTrend) String() string {
	return fmt.Sprintf(
		"UserWeeklyTrend{Id: %s, UserId: %s, WeekStart: %s, FocusMinutes: %.2f, MeetingMinutes: %.2f, BreakMinutes: %.2f, CreatedAt: %s, UpdatedAt: %s}",
		u.Id, u.UserId, u.WeekStart.Format(time.RFC3339), u.FocusMinutes, u.MeetingMinutes, u.BreakMinutes, u.CreatedAt.Format(time.RFC3339), u.UpdatedAt.Format(time.RFC3339),
	)
}

func (w WeeklyTrendDao) String() string {
	return fmt.Sprintf(
		"WeeklyTrendDao{UserId: %s, NoOfDays: %s}",
		w.UserId, w.NoOfDays,
	)
}
