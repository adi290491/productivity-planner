package models

type Repository interface {
	FetchDailyTrend(dailyTrendDao *DailyTrendDao) ([]UserDailyTrend, error)
	FetchWeeklyTrend(dailyTrendDao *WeeklyTrendDao) ([]UserWeeklyTrend, error)
}
