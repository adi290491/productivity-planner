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

/*
 Model object:
 [UserDailyTrend
 {Id: 0ee1d2fc-f645-4460-b881-743c1d1a78ed,
 UserId: b6ac7789-2453-47c2-b4d8-6371d07b4450,
  Day: 2025-05-20T00:00:00Z,
  FocusMinutes: 300.00,
   MeetingMinutes: 0.00,
   BreakMinutes: 150.00,
   CreatedAt: 2025-05-20T05:16:00Z,
    UpdatedAt: 2025-05-20T05:49:00Z
	}
	 UserDailyTrend{Id: 10e7244b-acce-49ee-ada5-524aa829ed69, UserId: b6ac7789-2453-47c2-b4d8-6371d07b4450, Day: 2025-05-21T00:00:00Z, FocusMinutes: 0.00, MeetingMinutes: 30.00, BreakMinutes: 0.00, CreatedAt: 2025-05-21T00:00:00Z, UpdatedAt: 2025-05-21T00:00:00Z}]
*/
