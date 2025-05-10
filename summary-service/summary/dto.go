package summary

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
