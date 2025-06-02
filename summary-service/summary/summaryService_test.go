package summary

import (
	models "productivity-planner/summary-service/model"
	"testing"
	"time"
)

func TestGetDailySessionSummary(t *testing.T) {
	svc := &SummaryService{
		Repo: &models.TestDBRepo{},
	}

	tests := []struct {
		name           string
		userId         string
		date           string
		expectError    bool
		expectNonEmpty bool
	}{
		{
			name:           "valid daily summary",
			userId:         "11111111-1111-1111-1111-111111111111",
			date:           time.Now().Format("2006-01-02"),
			expectError:    false,
			expectNonEmpty: true,
		},
		{
			name:        "invalid date format",
			userId:      "11111111-1111-1111-1111-111111111111",
			date:        "not-a-date",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := svc.GetDailySessionSummary(tc.userId, tc.date)
			if tc.expectError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.expectNonEmpty && result == nil {
				t.Errorf("expected non-nil result")
			}
		})
	}
}

func TestGetWeeklySessionSummary(t *testing.T) {
	svc := &SummaryService{
		Repo: &models.TestDBRepo{},
	}

	tests := []struct {
		name           string
		userId         string
		start          string
		expectError    bool
		expectNonEmpty bool
	}{
		{
			name:           "valid weekly summary",
			userId:         "11111111-1111-1111-1111-111111111111",
			start:          time.Now().Format("2006-01-02"),
			expectError:    false,
			expectNonEmpty: true,
		},
		{
			name:        "invalid start date format",
			userId:      "11111111-1111-1111-1111-111111111111",
			start:       "not-a-date",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := svc.GetWeeklySessionSummary(tc.userId, tc.start)
			if tc.expectError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.expectNonEmpty && result == nil {
				t.Errorf("expected non-nil result")
			}
		})
	}
}
