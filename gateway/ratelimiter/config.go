package ratelimiter

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Enabled         bool          // Whether rate limiting is enabled
	Capacity        float64       // Maximum tokens per bucket (burst size)
	RefillRate      float64       // Tokens added per second
	CleanupInterval time.Duration // How often to cleanup stale buckets
	TTL             time.Duration // Time before removing unused buckets
}

// LoadConfig loads rate limiter configuration from environment variables
// Returns default configuration if environment variables are not set
func LoadConfig() *Config {
	return &Config{
		Enabled:         getEnvBool("RATE_LIMIT_ENABLED", true),
		Capacity:        getEnvFloat("RATE_LIMIT_CAPACITY", 20.0),
		RefillRate:      getEnvFloat("RATE_LIMIT_REFILL_RATE", 10.0),
		CleanupInterval: getEnvDuration("RATE_LIMIT_CLEANUP_INTERVAL", 5*time.Minute),
		TTL:             getEnvDuration("RATE_LIMIT_TTL", 30*time.Minute),
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Capacity <= 0 {
		return fmt.Errorf("capacity must be positive, got %f", c.Capacity)
	}

	if c.RefillRate <= 0 {
		return fmt.Errorf("refill rate must be positive, got %f", c.RefillRate)
	}

	if c.CleanupInterval <= 0 {
		return fmt.Errorf("cleanup interval must be positive, got %v", c.CleanupInterval)
	}

	if c.TTL <= 0 {
		return fmt.Errorf("TTL must be positive, got %v", c.TTL)
	}

	if c.TTL < c.CleanupInterval {
		return fmt.Errorf("TTL (%v) must be >= cleanup interval (%v)", c.TTL, c.CleanupInterval)
	}

	return nil
}

// String returns the string representation of the config
func (c *Config) String() string {
	return fmt.Sprintf(
		"RateLimiterConfig{Enabled: %v, Capacity: %.1f, RefillRate: %.1f/s, CleanupInterval: %v, TTL: %v}",
		c.Enabled, c.Capacity, c.RefillRate, c.CleanupInterval, c.TTL,
	)
}

// getEnvBool gets a boolean from environment variable with a default value
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

// getEnvFloat gets a float from environment variable with a default value
func getEnvFloat(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return floatValue
}

// getEnvDuration gets a duration from environment variable with a default value
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return duration
}
