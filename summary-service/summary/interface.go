package summary

type SummaryServiceInterface interface {
	GetDailySessionSummary(userId string, date string) (*DailySessionSummary, error)
	GetWeeklySessionSummary(userId string, start string) (*WeeklySessionSummary, error)
}
