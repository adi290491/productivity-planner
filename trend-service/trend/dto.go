package trend

import "github.com/google/uuid"

type DailyTrendResponse struct {
	UserId uuid.UUID `json:"user_id"`
	Days   []Day     `json:"days"`
}

type Day struct {
	Date      string            `json:"date"`
	TotalTime string            `json:"total_time"`
	Breakdown map[string]string `json:"breakdown"`
}

type WeeklyTrendResponse struct {
	UserId uuid.UUID `json:"user_id"`
	Days   []Week    `json:"days"`
}

type Week struct {
	WeekStart string            `json:"week_start"`
	TotalTime string            `json:"total_time"`
	Breakdown map[string]string `json:"breakdown"`
}
