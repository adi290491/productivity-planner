package trend

import (
	"errors"
	"testing"

	models "github.com/adi290491/productivity-planner/trend-service/model"
	"github.com/google/uuid"
)

// Mock repository for testing
type MockRepository struct {
	dailyTrends            []models.UserDailyTrend
	weeklyTrends           []models.UserWeeklyTrend
	unviewedDailyCount     int
	unviewedWeeklyCount    int
	shouldReturnError      bool
	markDailyViewedCalled  bool
	markWeeklyViewedCalled bool
}

func (m *MockRepository) FetchDailyTrend(dailyTrendDao *models.DailyTrendDao) ([]models.UserDailyTrend, error) {
	if m.shouldReturnError {
		return nil, errors.New("database error")
	}
	return m.dailyTrends, nil
}

func (m *MockRepository) FetchWeeklyTrend(weeklyTrendDao *models.WeeklyTrendDao) ([]models.UserWeeklyTrend, error) {
	if m.shouldReturnError {
		return nil, errors.New("database error")
	}
	return m.weeklyTrends, nil
}

func (m *MockRepository) CountUnviewedDailyTrends(userID uuid.UUID) (int, error) {
	if m.shouldReturnError {
		return 0, errors.New("database error")
	}
	return m.unviewedDailyCount, nil
}

func (m *MockRepository) CountUnviewedWeeklyTrends(userID uuid.UUID) (int, error) {
	if m.shouldReturnError {
		return 0, errors.New("database error")
	}
	return m.unviewedWeeklyCount, nil
}

func (m *MockRepository) MarkDailyTrendsAsViewed(userId string) error {
	if m.shouldReturnError {
		return errors.New("database error")
	}
	m.markDailyViewedCalled = true
	return nil
}

func (m *MockRepository) MarkWeeklyTrendsAsViewed(userId string) error {
	if m.shouldReturnError {
		return errors.New("database error")
	}
	m.markWeeklyViewedCalled = true
	return nil
}

func TestTrendService_FetchDailyTrend(t *testing.T) {
	mockRepo := &MockRepository{
		dailyTrends: []models.UserDailyTrend{
			{
				UserId:         uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				FocusMinutes:   120,
				MeetingMinutes: 60,
				BreakMinutes:   30,
			},
		},
	}

	service := NewTrendService(mockRepo)

	t.Run("returns daily trends with valid parameters", func(t *testing.T) {
		result, err := service.FetchDailyTrend("11111111-1111-1111-1111-111111111111", "7")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result == nil {
			t.Error("expected result, got nil")
		}
		if result.UserId != "11111111-1111-1111-1111-111111111111" {
			t.Errorf("expected user id %s, got %s", "11111111-1111-1111-1111-111111111111", result.UserId)
		}
	})

	t.Run("returns error with invalid days parameter", func(t *testing.T) {
		_, err := service.FetchDailyTrend("11111111-1111-1111-1111-111111111111", "invalid")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		mockRepo.shouldReturnError = true
		_, err := service.FetchDailyTrend("11111111-1111-1111-1111-111111111111", "7")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTrendService_FetchWeeklyTrend(t *testing.T) {
	mockRepo := &MockRepository{
		weeklyTrends: []models.UserWeeklyTrend{
			{
				UserId:         uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				FocusMinutes:   600,
				MeetingMinutes: 300,
				BreakMinutes:   150,
			},
		},
	}

	service := NewTrendService(mockRepo)

	t.Run("returns weekly trends with valid parameters", func(t *testing.T) {
		result, err := service.FetchWeeklyTrend("11111111-1111-1111-1111-111111111111", "4")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result == nil {
			t.Error("expected result, got nil")
		}
		if result.UserId != "11111111-1111-1111-1111-111111111111" {
			t.Errorf("expected user id %s, got %s", "11111111-1111-1111-1111-111111111111", result.UserId)
		}
	})

	t.Run("returns error with invalid weeks parameter", func(t *testing.T) {
		_, err := service.FetchWeeklyTrend("11111111-1111-1111-1111-111111111111", "invalid")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		mockRepo.shouldReturnError = true
		_, err := service.FetchWeeklyTrend("11111111-1111-1111-1111-111111111111", "4")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestTrendService_GetUnviewedTrendsCount(t *testing.T) {
	mockRepo := &MockRepository{
		unviewedDailyCount:  5,
		unviewedWeeklyCount: 3,
	}

	service := NewTrendService(mockRepo)

	t.Run("returns unviewed counts with valid user id", func(t *testing.T) {
		result, err := service.GetUnviewedTrendsCount("11111111-1111-1111-1111-111111111111")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result == nil {
			t.Error("expected result, got nil")
		}
		if result.DailyCount != 5 {
			t.Errorf("expected daily count 5, got %d", result.DailyCount)
		}
		if result.WeeklyCount != 3 {
			t.Errorf("expected weekly count 3, got %d", result.WeeklyCount)
		}
	})

	t.Run("returns error when daily count fails", func(t *testing.T) {
		mockRepo.shouldReturnError = true
		_, err := service.GetUnviewedTrendsCount("11111111-1111-1111-1111-111111111111")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("returns error with invalid user id format", func(t *testing.T) {
		mockRepo.shouldReturnError = false
		_, err := service.GetUnviewedTrendsCount("invalid-uuid-format")

		if err == nil {
			t.Error("expected error for invalid UUID format, got nil")
		}
	})

	t.Run("returns error with empty user id", func(t *testing.T) {
		mockRepo.shouldReturnError = false
		_, err := service.GetUnviewedTrendsCount("")

		if err == nil {
			t.Error("expected error for empty UUID, got nil")
		}
	})
}

func TestTrendService_MarkTrendsAsViewed(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewTrendService(mockRepo)

	t.Run("marks daily trends as viewed", func(t *testing.T) {
		err := service.MarkTrendsAsViewed("11111111-1111-1111-1111-111111111111", "daily")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !mockRepo.markDailyViewedCalled {
			t.Error("expected MarkDailyTrendsAsViewed to be called")
		}
	})

	t.Run("marks weekly trends as viewed", func(t *testing.T) {
		mockRepo.markWeeklyViewedCalled = false
		err := service.MarkTrendsAsViewed("11111111-1111-1111-1111-111111111111", "weekly")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !mockRepo.markWeeklyViewedCalled {
			t.Error("expected MarkWeeklyTrendsAsViewed to be called")
		}
	})

	t.Run("returns error when repository fails for daily", func(t *testing.T) {
		mockRepo.shouldReturnError = true
		err := service.MarkTrendsAsViewed("11111111-1111-1111-1111-111111111111", "daily")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("returns error when repository fails for weekly", func(t *testing.T) {
		mockRepo.shouldReturnError = true
		err := service.MarkTrendsAsViewed("11111111-1111-1111-1111-111111111111", "weekly")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
