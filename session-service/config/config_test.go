package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_ValidEnvVars(t *testing.T) {
	os.Setenv("DB_HOSTNAME", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USERNAME", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("PORT", "1234")
	defer func() {
		os.Unsetenv("DB_HOSTNAME")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("PORT")
	}()

	cfg := Load()
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
}

func TestLoad_DefaultPortAndSSLMode(t *testing.T) {
	os.Setenv("DB_HOSTNAME", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USERNAME", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_SSLMODE", "false")
	os.Setenv("PORT", "")
	defer func() {
		os.Unsetenv("DB_HOSTNAME")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("PORT")
	}()

	cfg := Load()
	expectedDSN := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	if cfg.DSN != expectedDSN {
		t.Errorf("expected DSN %s, got %s", expectedDSN, cfg.DSN)
	}
	if cfg.Port != "8085" {
		t.Errorf("expected default port 8085, got %s", cfg.Port)
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
