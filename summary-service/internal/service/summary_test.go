package service

import (
	"context"
	"testing"
	"time"

	"github.com/adi290491/productivity-planner/summary-service/internal/model"
	"github.com/adi290491/productivity-planner/summary-service/internal/repository"
	"github.com/google/uuid"
)

func TestGetDailySummary_Success(t *testing.T) {
	endTime := time.Now().Add(-46 * time.Hour)
	mockRepo := &repository.MockRepository{
		Sessions: []model.Session{
			{
				ID:          uuid.New(),
				UserId:      uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				SessionType: "focus",
				StartTime:   time.Now().Add(-48 * time.Hour),
				EndTime:     &endTime,
			},
		},
	}

	svc := NewSummaryService(mockRepo)
	ctx := context.Background()

	result, err := svc.GetDailySessionSummary(ctx, "11111111-1111-1111-1111-111111111111", time.Now().Format("2006-01-02"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Date == "" {
		t.Error("expected non-empty date")
	}
}

func TestGetDailySummary_InvalidDate(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	svc := NewSummaryService(mockRepo)
	ctx := context.Background()

	_, err := svc.GetDailySessionSummary(ctx, "11111111-1111-1111-1111-111111111111", "invalid-date")
	if err == nil {
		t.Fatal("expected error for invalid date")
	}
}

func TestGetDailySummary_NoSessions(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	svc := NewSummaryService(mockRepo)
	ctx := context.Background()

	result, err := svc.GetDailySessionSummary(ctx, "notfound", time.Now().Format("2006-01-02"))
	if err == nil {
		t.Fatal("expected error for user with no sessions")
	}
	if result != nil {
		t.Error("expected nil result when error occurs")
	}
}

func TestGetWeeklySummary_Success(t *testing.T) {
	endTime := time.Now().Add(-46 * time.Hour)
	mockRepo := &repository.MockRepository{
		Sessions: []model.Session{
			{
				ID:          uuid.New(),
				UserId:      uuid.MustParse("11111111-1111-1111-1111-111111111111"),
				SessionType: "focus",
				StartTime:   time.Now().Add(-48 * time.Hour),
				EndTime:     &endTime,
			},
		},
	}

	svc := NewSummaryService(mockRepo)
	ctx := context.Background()

	result, err := svc.GetWeeklySessionSummary(ctx, "11111111-1111-1111-1111-111111111111", time.Now().Format("2006-01-02"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.StartDate == "" {
		t.Error("expected non-empty start date")
	}
	if result.EndDate == "" {
		t.Error("expected non-empty end date")
	}
}

func TestGetWeeklySummary_InvalidDate(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	svc := NewSummaryService(mockRepo)
	ctx := context.Background()

	_, err := svc.GetWeeklySessionSummary(ctx, "11111111-1111-1111-1111-111111111111", "invalid-date")
	if err == nil {
		t.Fatal("expected error for invalid date")
	}
}

func TestCalculateDailySummary(t *testing.T) {
	endTime1 := time.Now().Add(-1 * time.Hour)
	endTime2 := time.Now().Add(-30 * time.Minute)

	sessions := []model.Session{
		{
			ID:          uuid.New(),
			UserId:      uuid.New(),
			SessionType: "focus",
			StartTime:   time.Now().Add(-2 * time.Hour),
			EndTime:     &endTime1,
		},
		{
			ID:          uuid.New(),
			UserId:      uuid.New(),
			SessionType: "break",
			StartTime:   time.Now().Add(-45 * time.Minute),
			EndTime:     &endTime2,
		},
	}

	date := time.Now()
	result := calculateDailySummary(sessions, date)

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Date != date.Format("2006-01-02") {
		t.Errorf("expected date %s, got %s", date.Format("2006-01-02"), result.Date)
	}
	if len(result.Breakdown) != 2 {
		t.Errorf("expected 2 session types in breakdown, got %d", len(result.Breakdown))
	}
}

func TestCalculateWeeklySummary(t *testing.T) {
	endTime := time.Now().Add(-1 * time.Hour)
	sessions := []model.Session{
		{
			ID:          uuid.New(),
			UserId:      uuid.New(),
			SessionType: "focus",
			StartTime:   time.Now().Add(-2 * time.Hour),
			EndTime:     &endTime,
		},
	}

	start := time.Now().AddDate(0, 0, -3)
	end := time.Now().AddDate(0, 0, 4)

	result := calculateWeeklySummary(sessions, start, end)

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.StartDate != start.Format("2006-01-02") {
		t.Errorf("expected start date %s, got %s", start.Format("2006-01-02"), result.StartDate)
	}
	if result.EndDate != end.Format("2006-01-02") {
		t.Errorf("expected end date %s, got %s", end.Format("2006-01-02"), result.EndDate)
	}
}
