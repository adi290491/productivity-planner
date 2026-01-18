package ratelimiter

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear all rate limit env vars
	clearRateLimitEnvVars()

	config := LoadConfig()

	// Check defaults
	if !config.Enabled {
		t.Error("Expected Enabled to be true by default")
	}

	if config.Capacity != 20.0 {
		t.Errorf("Expected default Capacity 20.0, got %f", config.Capacity)
	}

	if config.RefillRate != 10.0 {
		t.Errorf("Expected default RefillRate 10.0, got %f", config.RefillRate)
	}

	if config.CleanupInterval != 5*time.Minute {
		t.Errorf("Expected default CleanupInterval 5m, got %v", config.CleanupInterval)
	}

	if config.TTL != 30*time.Minute {
		t.Errorf("Expected default TTL 30m, got %v", config.TTL)
	}
}

func TestLoadConfig_FromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("RATE_LIMIT_ENABLED", "false")
	os.Setenv("RATE_LIMIT_CAPACITY", "50")
	os.Setenv("RATE_LIMIT_REFILL_RATE", "25.5")
	os.Setenv("RATE_LIMIT_CLEANUP_INTERVAL", "10m")
	os.Setenv("RATE_LIMIT_TTL", "1h")

	defer clearRateLimitEnvVars()

	config := LoadConfig()

	if config.Enabled {
		t.Error("Expected Enabled to be false")
	}

	if config.Capacity != 50.0 {
		t.Errorf("Expected Capacity 50.0, got %f", config.Capacity)
	}

	if config.RefillRate != 25.5 {
		t.Errorf("Expected RefillRate 25.5, got %f", config.RefillRate)
	}

	if config.CleanupInterval != 10*time.Minute {
		t.Errorf("Expected CleanupInterval 10m, got %v", config.CleanupInterval)
	}

	if config.TTL != 1*time.Hour {
		t.Errorf("Expected TTL 1h, got %v", config.TTL)
	}
}

func TestLoadConfig_InvalidEnvValues(t *testing.T) {
	// Set invalid values - should fall back to defaults
	os.Setenv("RATE_LIMIT_ENABLED", "not-a-bool")
	os.Setenv("RATE_LIMIT_CAPACITY", "not-a-number")
	os.Setenv("RATE_LIMIT_REFILL_RATE", "invalid")
	os.Setenv("RATE_LIMIT_CLEANUP_INTERVAL", "invalid-duration")
	os.Setenv("RATE_LIMIT_TTL", "bad-format")

	defer clearRateLimitEnvVars()

	config := LoadConfig()

	// Should use defaults when parsing fails
	if !config.Enabled {
		t.Error("Expected Enabled to be true (default) when invalid")
	}

	if config.Capacity != 20.0 {
		t.Errorf("Expected Capacity 20.0 (default) when invalid, got %f", config.Capacity)
	}

	if config.RefillRate != 10.0 {
		t.Errorf("Expected RefillRate 10.0 (default) when invalid, got %f", config.RefillRate)
	}

	if config.CleanupInterval != 5*time.Minute {
		t.Errorf("Expected CleanupInterval 5m (default) when invalid, got %v", config.CleanupInterval)
	}

	if config.TTL != 30*time.Minute {
		t.Errorf("Expected TTL 30m (default) when invalid, got %v", config.TTL)
	}
}

func TestConfig_Validate_Valid(t *testing.T) {
	config := &Config{
		Enabled:         true,
		Capacity:        10.0,
		RefillRate:      5.0,
		CleanupInterval: 5 * time.Minute,
		TTL:             30 * time.Minute,
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Expected valid config to pass validation, got error: %v", err)
	}
}

func TestConfig_Validate_InvalidCapacity(t *testing.T) {
	config := &Config{
		Enabled:         true,
		Capacity:        0,
		RefillRate:      5.0,
		CleanupInterval: 5 * time.Minute,
		TTL:             30 * time.Minute,
	}

	err := config.Validate()
	if err == nil {
		t.Error("Expected validation error for zero capacity")
	}

	config.Capacity = -10.0
	err = config.Validate()
	if err == nil {
		t.Error("Expected validation error for negative capacity")
	}
}

func TestConfig_Validate_InvalidRefillRate(t *testing.T) {
	config := &Config{
		Enabled:         true,
		Capacity:        10.0,
		RefillRate:      0,
		CleanupInterval: 5 * time.Minute,
		TTL:             30 * time.Minute,
	}

	err := config.Validate()
	if err == nil {
		t.Error("Expected validation error for zero refill rate")
	}

	config.RefillRate = -5.0
	err = config.Validate()
	if err == nil {
		t.Error("Expected validation error for negative refill rate")
	}
}

func TestConfig_Validate_InvalidCleanupInterval(t *testing.T) {
	config := &Config{
		Enabled:         true,
		Capacity:        10.0,
		RefillRate:      5.0,
		CleanupInterval: 0,
		TTL:             30 * time.Minute,
	}

	err := config.Validate()
	if err == nil {
		t.Error("Expected validation error for zero cleanup interval")
	}

	config.CleanupInterval = -5 * time.Minute
	err = config.Validate()
	if err == nil {
		t.Error("Expected validation error for negative cleanup interval")
	}
}

func TestConfig_Validate_InvalidTTL(t *testing.T) {
	config := &Config{
		Enabled:         true,
		Capacity:        10.0,
		RefillRate:      5.0,
		CleanupInterval: 5 * time.Minute,
		TTL:             0,
	}

	err := config.Validate()
	if err == nil {
		t.Error("Expected validation error for zero TTL")
	}

	config.TTL = -30 * time.Minute
	err = config.Validate()
	if err == nil {
		t.Error("Expected validation error for negative TTL")
	}
}

func TestConfig_Validate_TTLLessThanCleanupInterval(t *testing.T) {
	config := &Config{
		Enabled:         true,
		Capacity:        10.0,
		RefillRate:      5.0,
		CleanupInterval: 30 * time.Minute,
		TTL:             5 * time.Minute, // TTL < CleanupInterval
	}

	err := config.Validate()
	if err == nil {
		t.Error("Expected validation error when TTL < CleanupInterval")
	}
}

func TestConfig_String(t *testing.T) {
	config := &Config{
		Enabled:         true,
		Capacity:        15.5,
		RefillRate:      7.2,
		CleanupInterval: 10 * time.Minute,
		TTL:             45 * time.Minute,
	}

	str := config.String()

	// Should contain key information
	if str == "" {
		t.Error("Expected non-empty string representation")
	}

	// Basic check that it contains the values
	t.Logf("Config string: %s", str)
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		{"empty uses default true", "", true, true},
		{"empty uses default false", "", false, false},
		{"true string", "true", false, true},
		{"false string", "false", true, false},
		{"1 as true", "1", false, true},
		{"0 as false", "0", true, false},
		{"invalid uses default", "invalid", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("TEST_BOOL", tt.envValue)
				defer os.Unsetenv("TEST_BOOL")
			} else {
				os.Unsetenv("TEST_BOOL")
			}

			result := getEnvBool("TEST_BOOL", tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetEnvFloat(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue float64
		expected     float64
	}{
		{"empty uses default", "", 10.0, 10.0},
		{"valid integer", "42", 10.0, 42.0},
		{"valid float", "3.14", 10.0, 3.14},
		{"negative float", "-5.5", 10.0, -5.5},
		{"invalid uses default", "not-a-number", 10.0, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("TEST_FLOAT", tt.envValue)
				defer os.Unsetenv("TEST_FLOAT")
			} else {
				os.Unsetenv("TEST_FLOAT")
			}

			result := getEnvFloat("TEST_FLOAT", tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestGetEnvDuration(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{"empty uses default", "", 5 * time.Minute, 5 * time.Minute},
		{"valid seconds", "30s", 5 * time.Minute, 30 * time.Second},
		{"valid minutes", "10m", 5 * time.Minute, 10 * time.Minute},
		{"valid hours", "2h", 5 * time.Minute, 2 * time.Hour},
		{"combined format", "1h30m", 5 * time.Minute, 90 * time.Minute},
		{"invalid uses default", "invalid", 5 * time.Minute, 5 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("TEST_DURATION", tt.envValue)
				defer os.Unsetenv("TEST_DURATION")
			} else {
				os.Unsetenv("TEST_DURATION")
			}

			result := getEnvDuration("TEST_DURATION", tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Helper function to clear all rate limit environment variables
func clearRateLimitEnvVars() {
	os.Unsetenv("RATE_LIMIT_ENABLED")
	os.Unsetenv("RATE_LIMIT_CAPACITY")
	os.Unsetenv("RATE_LIMIT_REFILL_RATE")
	os.Unsetenv("RATE_LIMIT_CLEANUP_INTERVAL")
	os.Unsetenv("RATE_LIMIT_TTL")
}
