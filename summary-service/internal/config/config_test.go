package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoad_ValidEnvVars(t *testing.T) {
	t.Setenv("DB_HOSTNAME", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("DB_USERNAME", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("PORT", "1234")
	t.Setenv("PROFILE", "test")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedDSN := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	if cfg.DSN != expectedDSN {
		t.Errorf("expected DSN %s, got %s", expectedDSN, cfg.DSN)
	}
	if cfg.Port != "1234" {
		t.Errorf("expected port 1234, got %s", cfg.Port)
	}
	if cfg.ReadTimeout != 10*time.Second {
		t.Errorf("expected ReadTimeout 10s, got %v", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 10*time.Second {
		t.Errorf("expected WriteTimeout 10s, got %v", cfg.WriteTimeout)
	}
	if cfg.Profile != "test" {
		t.Errorf("expected profile 'test', got %s", cfg.Profile)
	}
}

func TestLoad_DefaultPortAndSSLMode(t *testing.T) {
	t.Setenv("DB_HOSTNAME", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("DB_USERNAME", "testuser")
	t.Setenv("DB_PASSWORD", "testpass")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedDSN := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	if cfg.DSN != expectedDSN {
		t.Errorf("expected DSN %s, got %s", expectedDSN, cfg.DSN)
	}
	if cfg.Port != "8082" {
		t.Errorf("expected default port 8082, got %s", cfg.Port)
	}
	if cfg.Profile != "dev" {
		t.Errorf("expected default profile 'dev', got %s", cfg.Profile)
	}
}

func TestLoad_MissingRequiredEnvVars(t *testing.T) {
	// Clear all DB env vars
	for _, key := range []string{"DB_HOSTNAME", "DB_PORT", "DB_NAME", "DB_USERNAME", "DB_PASSWORD"} {
		os.Unsetenv(key)
	}

	cfg, err := Load()
	if cfg != nil {
		t.Errorf("expected nil config, got %+v", cfg)
	}
	if err == nil || !strings.Contains(err.Error(), "missing required database configuration") {
		t.Errorf("expected missing config error, got %v", err)
	}
}

func TestDBConfig_DSN(t *testing.T) {
	dbConf := &DBConfig{
		Host:     "localhost",
		Port:     "5432",
		DBName:   "testdb",
		User:     "testuser",
		Password: "testpass",
		SSLMode:  "disable",
	}
	expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	if dbConf.DSN() != expected {
		t.Errorf("expected %s, got %s", expected, dbConf.DSN())
	}
}

func TestRedactDSN(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid DSN with password",
			input:    "postgres://user:password123@localhost:5432/testdb",
			expected: "postgres://user:*****@localhost:5432/testdb",
		},
		{
			name:     "DSN without password",
			input:    "postgres://user@localhost:5432/testdb",
			expected: "postgres://user@localhost:5432/testdb",
		},
		{
			name:     "invalid DSN",
			input:    "not-a-valid-dsn",
			expected: "not-a-valid-dsn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := redactDSN(tt.input)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
