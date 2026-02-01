package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID
	UserId      uuid.UUID
	SessionType string
	StartTime   time.Time
	EndTime     *time.Time
}

type Summary struct {
	UserId    string
	StartTime time.Time
	EndTime   time.Time
}

// Stringer function for Session
func (s Session) String() string {
	return "Session{ID: " + s.ID.String() +
		", UserId: " + s.UserId.String() +
		", SessionType: " + s.SessionType +
		", StartTime: " + s.StartTime.String() +
		", EndTime: " + func() string {
		if s.EndTime != nil {
			return s.EndTime.String()
		}
		return "<nil>"
	}() + "}"
}

// Stringer function for Summary
func (s Summary) String() string {
	return "Summary{UserId: " + s.UserId +
		", StartTime: " + s.StartTime.String() +
		", EndTime: " + s.EndTime.String() + "}"
}

type DailySessionSummary struct {
	Date      string            `json:"date"`
	TotalTime string            `json:"total_time"`
	Breakdown map[string]string `json:"breakdown"`
}

type WeeklySessionSummary struct {
	StartDate      string                 `json:"start_date"`
	EndDate        string                 `json:"end_date"`
	TotalTime      string                 `json:"total_time"`
	DailySummaries []*DailySessionSummary `json:"daily_summaries"`
}
