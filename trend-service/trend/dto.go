package trend

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type DailyTrendResponse struct {
	UserId      string       `json:"user_id"`
	DailyTrends []DailyTrend `json:"dailyTrends"`
}

type DailyTrend struct {
	Date      string            `json:"date"`
	TotalTime string            `json:"total_time"`
	Breakdown map[string]string `json:"breakdown"`
}

type WeeklyTrendResponse struct {
	UserId       uuid.UUID     `json:"user_id"`
	WeeklyTrends []WeeklyTrend `json:"weekly_trends"`
}

type WeeklyTrend struct {
	WeekStart string            `json:"week_start"`
	TotalTime string            `json:"total_time"`
	Breakdown map[string]string `json:"breakdown"`
}

func (d DailyTrendResponse) String() string {
	var days []string
	for _, day := range d.DailyTrends {
		days = append(days, day.String())
	}
	return fmt.Sprintf("DailyTrendResponse{UserId: %s, Days: [%s]}", d.UserId, strings.Join(days, ", "))
}

func (d DailyTrend) String() string {
	return fmt.Sprintf("Day{Date: %s, TotalTime: %s, Breakdown: %v}", d.Date, d.TotalTime, d.Breakdown)
}

func (w WeeklyTrendResponse) String() string {
	var weeks []string
	for _, week := range w.WeeklyTrends {
		weeks = append(weeks, week.String())
	}
	return fmt.Sprintf("WeeklyTrendResponse{UserId: %s, Days: [%s]}", w.UserId, strings.Join(weeks, ", "))
}

func (w WeeklyTrend) String() string {
	return fmt.Sprintf("Week{WeekStart: %s, TotalTime: %s, Breakdown: %v}", w.WeekStart, w.TotalTime, w.Breakdown)
}
