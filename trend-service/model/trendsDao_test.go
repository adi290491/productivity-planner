package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTestDBRepo_FetchDailyTrend_ReturnsExpectedData(t *testing.T) {
	repo := &TestDBRepo{}

	dao := &DailyTrendDao{}
	result, err := repo.FetchDailyTrend(dao)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	expectedUUID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	if result[0].UserId != expectedUUID {
		t.Errorf("expected UserId %v, got %v", expectedUUID, result[0].UserId)
	}

	// Check that the Day is within the last 2 days (since it uses time.Now().AddDate(0,0,-1))
	now := time.Now()
	day := result[0].Day
	if day.Before(now.AddDate(0, 0, -2)) || day.After(now) {
		t.Errorf("expected Day to be within the last 2 days, got %v", day)
	}

	if result[0].FocusMinutes != 60 {
		t.Errorf("expected FocusMinutes 60, got %v", result[0].FocusMinutes)
	}
	if result[0].MeetingMinutes != 30 {
		t.Errorf("expected MeetingMinutes 30, got %v", result[0].MeetingMinutes)
	}
	if result[0].BreakMinutes != 15 {
		t.Errorf("expected BreakMinutes 15, got %v", result[0].BreakMinutes)
	}
}

func TestTestDBRepo_FetchDailyTrend_NilDao(t *testing.T) {
	repo := &TestDBRepo{}
	// Passing nil should not cause a panic or error
	result, err := repo.FetchDailyTrend(nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty, got %d", len(result))
	}
}
