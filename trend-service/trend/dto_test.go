package trend

import (
	"strings"
	"testing"
)

func TestDailyTrendResponse_String(t *testing.T) {
	t.Run("formats response with multiple trends", func(t *testing.T) {
		response := DailyTrendResponse{
			UserId: "test-user",
			DailyTrends: []DailyTrend{
				{Date: "2024-01-15", TotalTime: "3h 30m", Breakdown: map[string]string{"focus": "2h", "meeting": "1h", "break": "30m"}},
				{Date: "2024-01-16", TotalTime: "4h 0m", Breakdown: map[string]string{"focus": "3h", "meeting": "1h", "break": "0m"}},
			},
		}

		result := response.String()

		if !strings.Contains(result, "test-user") {
			t.Error("expected result to contain user id")
		}
		if !strings.Contains(result, "DailyTrendResponse") {
			t.Error("expected result to contain type name")
		}
	})

	t.Run("formats response with empty trends", func(t *testing.T) {
		response := DailyTrendResponse{
			UserId:      "test-user",
			DailyTrends: []DailyTrend{},
		}

		result := response.String()

		if !strings.Contains(result, "test-user") {
			t.Error("expected result to contain user id")
		}
		if !strings.Contains(result, "DailyTrendResponse") {
			t.Error("expected result to contain type name")
		}
	})
}

func TestDailyTrend_String(t *testing.T) {
	t.Run("formats trend correctly", func(t *testing.T) {
		trend := DailyTrend{
			Date:      "2024-01-15",
			TotalTime: "3h 30m",
			Breakdown: map[string]string{
				"focus":   "2h",
				"meeting": "1h",
				"break":   "30m",
			},
		}

		result := trend.String()

		if !strings.Contains(result, "2024-01-15") {
			t.Error("expected result to contain date")
		}
		if !strings.Contains(result, "3h 30m") {
			t.Error("expected result to contain total time")
		}
		if !strings.Contains(result, "Day{") {
			t.Error("expected result to contain type identifier")
		}
	})
}

func TestWeeklyTrendResponse_String(t *testing.T) {
	t.Run("formats response with multiple trends", func(t *testing.T) {
		response := WeeklyTrendResponse{
			UserId: "test-user",
			WeeklyTrends: []WeeklyTrend{
				{WeekStart: "2024-01-15", TotalTime: "20h 30m", Breakdown: map[string]string{"focus": "15h", "meeting": "5h", "break": "30m"}},
				{WeekStart: "2024-01-22", TotalTime: "25h 0m", Breakdown: map[string]string{"focus": "20h", "meeting": "5h", "break": "0m"}},
			},
		}

		result := response.String()

		if !strings.Contains(result, "test-user") {
			t.Error("expected result to contain user id")
		}
		if !strings.Contains(result, "WeeklyTrendResponse") {
			t.Error("expected result to contain type name")
		}
	})

	t.Run("formats response with empty trends", func(t *testing.T) {
		response := WeeklyTrendResponse{
			UserId:       "test-user",
			WeeklyTrends: []WeeklyTrend{},
		}

		result := response.String()

		if !strings.Contains(result, "test-user") {
			t.Error("expected result to contain user id")
		}
		if !strings.Contains(result, "WeeklyTrendResponse") {
			t.Error("expected result to contain type name")
		}
	})
}

func TestWeeklyTrend_String(t *testing.T) {
	t.Run("formats trend correctly", func(t *testing.T) {
		trend := WeeklyTrend{
			WeekStart: "2024-01-15",
			TotalTime: "20h 30m",
			Breakdown: map[string]string{
				"focus":   "15h",
				"meeting": "5h",
				"break":   "30m",
			},
		}

		result := trend.String()

		if !strings.Contains(result, "2024-01-15") {
			t.Error("expected result to contain week start")
		}
		if !strings.Contains(result, "20h 30m") {
			t.Error("expected result to contain total time")
		}
		if !strings.Contains(result, "Week{") {
			t.Error("expected result to contain type identifier")
		}
	})
}

func TestUnviewedTrendsCount(t *testing.T) {
	t.Run("creates count struct correctly", func(t *testing.T) {
		count := UnviewedTrendsCount{
			DailyCount:  5,
			WeeklyCount: 3,
		}

		if count.DailyCount != 5 {
			t.Errorf("expected daily count 5, got %d", count.DailyCount)
		}
		if count.WeeklyCount != 3 {
			t.Errorf("expected weekly count 3, got %d", count.WeeklyCount)
		}
	})

	t.Run("handles zero counts", func(t *testing.T) {
		count := UnviewedTrendsCount{
			DailyCount:  0,
			WeeklyCount: 0,
		}

		if count.DailyCount != 0 {
			t.Errorf("expected daily count 0, got %d", count.DailyCount)
		}
		if count.WeeklyCount != 0 {
			t.Errorf("expected weekly count 0, got %d", count.WeeklyCount)
		}
	})
}
