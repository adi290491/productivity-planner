package models

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUserDailyTrend_String(t *testing.T) {
	t.Run("formats daily trend correctly", func(t *testing.T) {
		id := uuid.New()
		userId := uuid.New()
		day := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
		createdAt := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)

		trend := UserDailyTrend{
			Id:             id,
			UserId:         userId,
			Day:            day,
			FocusMinutes:   120.50,
			MeetingMinutes: 60.25,
			BreakMinutes:   30.75,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
		}

		result := trend.String()

		if !strings.Contains(result, "UserDailyTrend{") {
			t.Error("expected result to contain type identifier")
		}
		if !strings.Contains(result, id.String()) {
			t.Error("expected result to contain ID")
		}
		if !strings.Contains(result, userId.String()) {
			t.Error("expected result to contain UserId")
		}
		if !strings.Contains(result, "2024-01-15T10:30:00Z") {
			t.Error("expected result to contain formatted day")
		}
		if !strings.Contains(result, "120.50") {
			t.Error("expected result to contain focus minutes")
		}
		if !strings.Contains(result, "60.25") {
			t.Error("expected result to contain meeting minutes")
		}
		if !strings.Contains(result, "30.75") {
			t.Error("expected result to contain break minutes")
		}
	})
}

func TestDailyTrendDao_String(t *testing.T) {
	t.Run("formats daily trend DAO correctly", func(t *testing.T) {
		lookbackDays := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
		dao := DailyTrendDao{
			UserId:       "test-user-123",
			LookbackDays: lookbackDays,
		}

		result := dao.String()

		if !strings.Contains(result, "DailyTrendDao{") {
			t.Error("expected result to contain type identifier")
		}
		if !strings.Contains(result, "test-user-123") {
			t.Error("expected result to contain UserId")
		}
		if !strings.Contains(result, "2024-01-10T00:00:00Z") {
			t.Error("expected result to contain formatted lookback days")
		}
	})
}

func TestUserWeeklyTrend_String(t *testing.T) {
	t.Run("formats weekly trend correctly", func(t *testing.T) {
		id := uuid.New()
		userId := uuid.New()
		weekStart := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		createdAt := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)

		trend := UserWeeklyTrend{
			Id:             id,
			UserId:         userId,
			WeekStart:      weekStart,
			FocusMinutes:   600.75,
			MeetingMinutes: 300.50,
			BreakMinutes:   150.25,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
		}

		result := trend.String()

		if !strings.Contains(result, "UserWeeklyTrend{") {
			t.Error("expected result to contain type identifier")
		}
		if !strings.Contains(result, id.String()) {
			t.Error("expected result to contain ID")
		}
		if !strings.Contains(result, userId.String()) {
			t.Error("expected result to contain UserId")
		}
		if !strings.Contains(result, "2024-01-15T00:00:00Z") {
			t.Error("expected result to contain formatted week start")
		}
		if !strings.Contains(result, "600.75") {
			t.Error("expected result to contain focus minutes")
		}
		if !strings.Contains(result, "300.50") {
			t.Error("expected result to contain meeting minutes")
		}
		if !strings.Contains(result, "150.25") {
			t.Error("expected result to contain break minutes")
		}
	})
}

func TestWeeklyTrendDao_String(t *testing.T) {
	t.Run("formats weekly trend DAO correctly", func(t *testing.T) {
		lookbackWeeks := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		dao := WeeklyTrendDao{
			UserId:        "test-user-456",
			LookbackWeeks: lookbackWeeks,
		}

		result := dao.String()

		if !strings.Contains(result, "WeeklyTrendDao{") {
			t.Error("expected result to contain type identifier")
		}
		if !strings.Contains(result, "test-user-456") {
			t.Error("expected result to contain UserId")
		}
		if !strings.Contains(result, "2024-01-01T00:00:00Z") {
			t.Error("expected result to contain formatted lookback weeks")
		}
	})
}

func TestUserDailyTrend_ViewedAtField(t *testing.T) {
	t.Run("handles nil ViewedAt field", func(t *testing.T) {
		trend := UserDailyTrend{
			Id:       uuid.New(),
			UserId:   uuid.New(),
			Day:      time.Now(),
			ViewedAt: nil, // This should not cause panic
		}

		// Test that the struct can be created and accessed without panicking
		if trend.ViewedAt != nil {
			t.Error("expected ViewedAt to be nil")
		}
	})

	t.Run("handles set ViewedAt field", func(t *testing.T) {
		now := time.Now()
		trend := UserDailyTrend{
			Id:       uuid.New(),
			UserId:   uuid.New(),
			Day:      time.Now(),
			ViewedAt: &now,
		}

		if trend.ViewedAt == nil {
			t.Error("expected ViewedAt to be set")
		}
		if !trend.ViewedAt.Equal(now) {
			t.Error("expected ViewedAt to match set time")
		}
	})
}

func TestUserWeeklyTrend_ViewedAtField(t *testing.T) {
	t.Run("handles nil ViewedAt field", func(t *testing.T) {
		trend := UserWeeklyTrend{
			Id:        uuid.New(),
			UserId:    uuid.New(),
			WeekStart: time.Now(),
			ViewedAt:  nil, // This should not cause panic
		}

		// Test that the struct can be created and accessed without panicking
		if trend.ViewedAt != nil {
			t.Error("expected ViewedAt to be nil")
		}
	})

	t.Run("handles set ViewedAt field", func(t *testing.T) {
		now := time.Now()
		trend := UserWeeklyTrend{
			Id:        uuid.New(),
			UserId:    uuid.New(),
			WeekStart: time.Now(),
			ViewedAt:  &now,
		}

		if trend.ViewedAt == nil {
			t.Error("expected ViewedAt to be set")
		}
		if !trend.ViewedAt.Equal(now) {
			t.Error("expected ViewedAt to match set time")
		}
	})
}
