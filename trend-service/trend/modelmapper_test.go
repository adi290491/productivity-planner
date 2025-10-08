package trend

import (
	"testing"
	"time"

	models "github.com/adi290491/productivity-planner/trend-service/model"
	"github.com/google/uuid"
)

func TestMapModelToResponse(t *testing.T) {
	userId := "11111111-1111-1111-1111-111111111111"

	t.Run("returns empty response when no trends", func(t *testing.T) {
		var emptyTrends []models.UserDailyTrend
		response := MapModelToResponse(emptyTrends, userId)

		if response.UserId != userId {
			t.Errorf("expected userId %s, got %s", userId, response.UserId)
		}
		if len(response.DailyTrends) != 0 {
			t.Errorf("expected empty trends, got %d trends", len(response.DailyTrends))
		}
	})

	t.Run("maps single trend correctly", func(t *testing.T) {
		trends := []models.UserDailyTrend{
			{
				Id:             uuid.New(),
				UserId:         uuid.MustParse(userId),
				Day:            time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				FocusMinutes:   120, // 2 hours
				MeetingMinutes: 60,  // 1 hour
				BreakMinutes:   30,  // 30 minutes
			},
		}

		response := MapModelToResponse(trends, userId)

		if response.UserId != userId {
			t.Errorf("expected userId %s, got %s", userId, response.UserId)
		}
		if len(response.DailyTrends) != 1 {
			t.Errorf("expected 1 trend, got %d trends", len(response.DailyTrends))
		}

		trend := response.DailyTrends[0]
		if trend.Date != "2024-01-15" {
			t.Errorf("expected date 2024-01-15, got %s", trend.Date)
		}
		if trend.Breakdown["focus"] == "" {
			t.Error("expected focus time to be set")
		}
		if trend.Breakdown["meeting"] == "" {
			t.Error("expected meeting time to be set")
		}
		if trend.Breakdown["break"] == "" {
			t.Error("expected break time to be set")
		}
	})

	t.Run("sorts trends by date in ascending order", func(t *testing.T) {
		trends := []models.UserDailyTrend{
			{
				Id:     uuid.New(),
				UserId: uuid.MustParse(userId),
				Day:    time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC),
			},
			{
				Id:     uuid.New(),
				UserId: uuid.MustParse(userId),
				Day:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			{
				Id:     uuid.New(),
				UserId: uuid.MustParse(userId),
				Day:    time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
			},
		}

		response := MapModelToResponse(trends, userId)

		if len(response.DailyTrends) != 3 {
			t.Errorf("expected 3 trends, got %d trends", len(response.DailyTrends))
		}

		expectedDates := []string{"2024-01-15", "2024-01-16", "2024-01-17"}
		for i, trend := range response.DailyTrends {
			if trend.Date != expectedDates[i] {
				t.Errorf("expected date %s at index %d, got %s", expectedDates[i], i, trend.Date)
			}
		}
	})
}

func TestMapWeeklyModelToResponse(t *testing.T) {
	userId := "11111111-1111-1111-1111-111111111111"

	t.Run("returns empty response when no trends", func(t *testing.T) {
		var emptyTrends []models.UserWeeklyTrend
		response := MapWeeklyModelToResponse(emptyTrends, userId)

		if response.UserId != userId {
			t.Errorf("expected userId %s, got %s", userId, response.UserId)
		}
		if len(response.WeeklyTrends) != 0 {
			t.Errorf("expected empty trends, got %d trends", len(response.WeeklyTrends))
		}
	})

	t.Run("maps single weekly trend correctly", func(t *testing.T) {
		trends := []models.UserWeeklyTrend{
			{
				Id:             uuid.New(),
				UserId:         uuid.MustParse(userId),
				WeekStart:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				FocusMinutes:   600, // 10 hours
				MeetingMinutes: 300, // 5 hours
				BreakMinutes:   150, // 2.5 hours
			},
		}

		response := MapWeeklyModelToResponse(trends, userId)

		if response.UserId != userId {
			t.Errorf("expected userId %s, got %s", userId, response.UserId)
		}
		if len(response.WeeklyTrends) != 1 {
			t.Errorf("expected 1 trend, got %d trends", len(response.WeeklyTrends))
		}

		trend := response.WeeklyTrends[0]
		if trend.WeekStart != "2024-01-15" {
			t.Errorf("expected week start 2024-01-15, got %s", trend.WeekStart)
		}
		if trend.Breakdown["focus"] == "" {
			t.Error("expected focus time to be set")
		}
		if trend.Breakdown["meeting"] == "" {
			t.Error("expected meeting time to be set")
		}
		if trend.Breakdown["break"] == "" {
			t.Error("expected break time to be set")
		}
	})

	t.Run("sorts weekly trends by week start in ascending order", func(t *testing.T) {
		trends := []models.UserWeeklyTrend{
			{
				Id:        uuid.New(),
				UserId:    uuid.MustParse(userId),
				WeekStart: time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC),
			},
			{
				Id:        uuid.New(),
				UserId:    uuid.MustParse(userId),
				WeekStart: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			{
				Id:        uuid.New(),
				UserId:    uuid.MustParse(userId),
				WeekStart: time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC),
			},
		}

		response := MapWeeklyModelToResponse(trends, userId)

		if len(response.WeeklyTrends) != 3 {
			t.Errorf("expected 3 trends, got %d trends", len(response.WeeklyTrends))
		}

		expectedWeeks := []string{"2024-01-15", "2024-01-22", "2024-01-29"}
		for i, trend := range response.WeeklyTrends {
			if trend.WeekStart != expectedWeeks[i] {
				t.Errorf("expected week start %s at index %d, got %s", expectedWeeks[i], i, trend.WeekStart)
			}
		}
	})
}
