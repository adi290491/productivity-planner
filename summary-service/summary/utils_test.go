package summary

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