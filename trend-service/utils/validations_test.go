package utils

import (
	"math"
	"testing"
)

func TestValidateDays_ValidValues(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1", 1},
		{"15", 15},
		{"30", 30},
	}
	for _, tt := range tests {
		got, err := ValidateDays(tt.input)
		if err != nil {
			t.Errorf("ValidateDays(%q) unexpected error: %v", tt.input, err)
		}
		if got != tt.expected {
			t.Errorf("ValidateDays(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestValidateDays_InvalidValues(t *testing.T) {
	tests := []string{
		"0",   // below min
		"31",  // above max
		"-5",  // negative
		"abc", // not a number
		"",    // empty
		"1.5", // float
	}
	for _, input := range tests {
		got, err := ValidateDays(input)
		if err == nil {
			t.Errorf("ValidateDays(%q) expected error, got nil", input)
		}
		if !math.IsNaN(float64(got)) && got != int(math.NaN()) {
			t.Errorf("ValidateDays(%q) expected NaN, got %d", input, got)
		}
	}
}

func TestValidateWeeks_ValidValues(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1", 1},
		{"6", 6},
		{"12", 12},
	}
	for _, tt := range tests {
		got, err := ValidateWeeks(tt.input)
		if err != nil {
			t.Errorf("ValidateWeeks(%q) unexpected error: %v", tt.input, err)
		}
		if got != tt.expected {
			t.Errorf("ValidateWeeks(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestValidateWeeks_InvalidValues(t *testing.T) {
	tests := []string{
		"0",   // below min
		"13",  // above max
		"-2",  // negative
		"xyz", // not a number
		"",    // empty
		"2.5", // float
	}
	for _, input := range tests {
		got, err := ValidateWeeks(input)
		if err == nil {
			t.Errorf("ValidateWeeks(%q) expected error, got nil", input)
		}
		if !math.IsNaN(float64(got)) && got != int(math.NaN()) {
			t.Errorf("ValidateWeeks(%q) expected NaN, got %d", input, got)
		}
	}
}

func TestFormatTimeToHrMin(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "00h00m"},
		{5, "00h05m"},
		{60, "01h00m"},
		{75, "01h15m"},
		{135, "02h15m"},
		{1439, "23h59m"},
	}
	for _, tt := range tests {
		got := FormatTimeToHrMin(tt.input)
		if got != tt.expected {
			t.Errorf("FormatTimeToHrMin(%v) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
