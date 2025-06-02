package summary

import (
	"fmt"
	models "productivity-planner/summary-service/model"
)

type MockSummaryService struct {
	Repo models.Repository
}

func (s *MockSummaryService) GetDailySessionSummary(userId string, date string) (*DailySessionSummary, error) {

	if date == "invalid-date" {
		return nil, fmt.Errorf("invalid date format")
	}

	if userId == "notfound" {
		return nil, fmt.Errorf("no sessions found for the given day")
	}

	return &DailySessionSummary{
		Date:      date,
		TotalTime: "2h30m",
		Breakdown: map[string]string{
			"focus":   "1h30m",
			"meeting": "1h0m",
		},
	}, nil
}

func (s *MockSummaryService) GetWeeklySessionSummary(userId string, start string) (*WeeklySessionSummary, error) {

	if start == "invalid-date" {
		return nil, fmt.Errorf("invalid date format")
	}

	if userId == "notfound" {
		return nil, fmt.Errorf("no sessions found")
	}

	return &WeeklySessionSummary{
		StartDate: "2025-05-19",
		EndDate:   "2025-05-25",
		TotalTime: "10h0m",
		DailySummaries: []*DailySessionSummary{
			{
				Date:      "2025-05-19",
				TotalTime: "2h0m",
				Breakdown: map[string]string{
					"focus": "1h0m",
				},
			},
		},
	}, nil
}
