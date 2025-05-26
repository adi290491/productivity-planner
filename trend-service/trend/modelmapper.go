package trend

import (
	"log"
	models "productivity-planner/trend-service/model"
	"productivity-planner/trend-service/utils"
	"sort"
)

func MapModelToResponse(dailyUserTrend []models.UserDailyTrend, userId string) *DailyTrendResponse {

	log.Println("Model object:", dailyUserTrend)

	if len(dailyUserTrend) == 0 {
		return &DailyTrendResponse{
			UserId:      userId,
			DailyTrends: []DailyTrend{},
		}
	}

	dailyTrends := make([]DailyTrend, 0)

	for _, trend := range dailyUserTrend {
		dailyTrends = append(dailyTrends, DailyTrend{
			Date:      trend.Day.Format("2006-01-02"),
			TotalTime: utils.FormatTimeToHrMin(trend.FocusMinutes + trend.MeetingMinutes + trend.BreakMinutes),
			Breakdown: map[string]string{
				"focus":   utils.FormatTimeToHrMin(trend.FocusMinutes),
				"meeting": utils.FormatTimeToHrMin(trend.MeetingMinutes),
				"break":   utils.FormatTimeToHrMin(trend.BreakMinutes),
			},
		})
	}
	// Sort the daily trends by date in ascending order
	sort.Slice(dailyTrends, func(i, j int) bool {
		return dailyTrends[i].Date < dailyTrends[j].Date
	})

	dailyTrendResponse := &DailyTrendResponse{
		UserId:      userId,
		DailyTrends: dailyTrends,
	}
	return dailyTrendResponse
}

func MapWeeklyModelToResponse(weeklyUserTrend []models.UserWeeklyTrend, userId string) *WeeklyTrendResponse {

	log.Println("Model object:", weeklyUserTrend)

	if len(weeklyUserTrend) == 0 {
		return &WeeklyTrendResponse{
			UserId:       userId,
			WeeklyTrends: []WeeklyTrend{},
		}
	}

	weeklyTrends := make([]WeeklyTrend, 0)
	for _, trend := range weeklyUserTrend {
		weeklyTrends = append(weeklyTrends, WeeklyTrend{
			WeekStart: trend.WeekStart.Format("2006-01-02"),
			TotalTime: utils.FormatTimeToHrMin(trend.FocusMinutes + trend.MeetingMinutes + trend.BreakMinutes),
			Breakdown: map[string]string{
				"focus":   utils.FormatTimeToHrMin(trend.FocusMinutes),
				"meeting": utils.FormatTimeToHrMin(trend.MeetingMinutes),
				"break":   utils.FormatTimeToHrMin(trend.BreakMinutes),
			},
		})
	}
	// Sort the weekly trends by week start date in ascending order
	sort.Slice(weeklyTrends, func(i, j int) bool {
		return weeklyTrends[i].WeekStart < weeklyTrends[j].WeekStart
	})

	weeklyTrendsResponse := &WeeklyTrendResponse{
		UserId:       userId,
		WeeklyTrends: weeklyTrends,
	}
	return weeklyTrendsResponse
}
