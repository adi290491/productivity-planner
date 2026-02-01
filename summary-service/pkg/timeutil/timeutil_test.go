package timeutil

import (
	"testing"
	"time"
)

func TestStartOfDayUTC_ValidDate(t *testing.T) {
	dateStr := "2024-06-10"
	got, err := StartOfDayUTC(dateStr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestStartOfDayUTC_EmptyDate(t *testing.T) {
	got, err := StartOfDayUTC("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	now := time.Now().UTC()
	expected := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if got.Year() != expected.Year() || got.Month() != expected.Month() || got.Day() != expected.Day() ||
		got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestStartOfDayUTC_InvalidDate(t *testing.T) {
	_, err := StartOfDayUTC("invalid-date")
	if err == nil {
		t.Errorf("expected error for invalid date, got nil")
	}
}

func TestEndOfDayUTC(t *testing.T) {
	input := time.Date(2024, 6, 10, 12, 30, 45, 0, time.UTC)
	got := EndOfDayUTC(input)
	expected := time.Date(2024, 6, 10, 23, 59, 59, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestStartOfWeekUTC_ValidDate(t *testing.T) {
	// June 13, 2024 is a Thursday
	dateStr := "2024-06-13"
	got, err := StartOfWeekUTC(dateStr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should return Monday, June 10, 2024
	expected := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestStartOfWeekUTC_Monday(t *testing.T) {
	// June 10, 2024 is a Monday
	dateStr := "2024-06-10"
	got, err := StartOfWeekUTC(dateStr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should return the same Monday
	expected := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestStartOfWeekUTC_Sunday(t *testing.T) {
	// June 16, 2024 is a Sunday
	dateStr := "2024-06-16"
	got, err := StartOfWeekUTC(dateStr)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should return Monday, June 10, 2024
	expected := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestStartOfWeekUTC_InvalidDate(t *testing.T) {
	_, err := StartOfWeekUTC("not-a-date")
	if err == nil {
		t.Errorf("expected error for invalid date, got nil")
	}
}

func TestDurationToHumanFormat(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "2 hours 30 minutes",
			duration: 2*time.Hour + 30*time.Minute,
			expected: "2h30m",
		},
		{
			name:     "1 hour",
			duration: 1 * time.Hour,
			expected: "1h0m",
		},
		{
			name:     "45 minutes",
			duration: 45 * time.Minute,
			expected: "0h45m",
		},
		{
			name:     "0 duration",
			duration: 0,
			expected: "0h0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DurationToHumanFormat(tt.duration)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
