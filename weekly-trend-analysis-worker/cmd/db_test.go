package main

import (
	"net/url"
	"os"
	"strings"
	"testing"

	"gorm.io/gorm"
)

func setEnvVars(vars map[string]string) func() {
	orig := make(map[string]string)
	for k, v := range vars {
		orig[k] = os.Getenv(k)
		os.Setenv(k, v)
	}
	return func() {
		for k, v := range orig {
			os.Setenv(k, v)
		}
	}
}

func unsetEnvVars(keys []string) func() {
	orig := make(map[string]string)
	for _, k := range keys {
		orig[k] = os.Getenv(k)
		os.Unsetenv(k)
	}
	return func() {
		for k, v := range orig {
			os.Setenv(k, v)
		}
	}
}

func TestDBConfig_DSN(t *testing.T) {
	cfg := &DBConfig{
		Host:     "localhost",
		Port:     "5432",
		DBName:   "testdb",
		User:     "testuser",
		Password: "testpass",
		SSLMode:  "disable",
	}
	dsn := cfg.DSN()
	expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	if dsn != expected {
		t.Errorf("DSN not constructed properly: got %s, want %s", dsn, expected)
	}
}

func TestRedactDSN(t *testing.T) {
	dsn := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	redacted := redactDSN(dsn)
	if !strings.Contains(redacted, url.QueryEscape("*****")) {
		t.Errorf("DSN not redacted properly: %s", redacted)
	}
	// Should not redact if no password
	dsnNoPass := "postgres://testuser@localhost:5432/testdb?sslmode=disable"
	if redactDSN(dsnNoPass) != dsnNoPass {
		t.Errorf("DSN without password should not be changed")
	}
}

func TestLoadConfig_AllEnvSet(t *testing.T) {
	cleanup := setEnvVars(map[string]string{
		"DB_HOSTNAME": "localhost",
		"DB_PORT":     "5432",
		"DB_NAME":     "testdb",
		"DB_USERNAME": "testuser",
		"DB_PASSWORD": "testpass",
		"DB_SSLMODE":  "require",
		"PROFILE":     "test",
	})
	defer cleanup()

	app, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app == nil {
		t.Fatal("expected app config, got nil")
	}
	if !strings.Contains(app.DSN, "require") {
		t.Errorf("expected sslmode=require in DSN, got %s", app.DSN)
	}
}

func TestLoadConfig_DefaultSSLMode(t *testing.T) {
	cleanup := setEnvVars(map[string]string{
		"DB_HOSTNAME": "localhost",
		"DB_PORT":     "5432",
		"DB_NAME":     "testdb",
		"DB_USERNAME": "testuser",
		"DB_PASSWORD": "testpass",
		"DB_SSLMODE":  "",
	})
	defer cleanup()

	app, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(app.DSN, "sslmode=disable") {
		t.Errorf("expected default sslmode=disable, got %s", app.DSN)
	}
}

func TestLoadConfig_MissingRequiredEnv(t *testing.T) {
	keys := []string{"DB_HOSTNAME", "DB_PORT", "DB_NAME", "DB_USERNAME", "DB_PASSWORD"}
	cleanup := unsetEnvVars(keys)
	defer cleanup()

	app, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for missing required envs, got nil")
	}
	if app != nil {
		t.Errorf("expected nil app on error, got: %v", app)
	}
}

func TestApplication_StructFields(t *testing.T) {
	app := &Application{DSN: "dsn-string", DB: nil}
	if app.DSN != "dsn-string" {
		t.Errorf("expected DSN 'dsn-string', got %s", app.DSN)
	}
	if app.DB != nil {
		t.Errorf("expected DB to be nil, got %v", app.DB)
	}
}

func TestApplication_String(t *testing.T) {
	app := &Application{DSN: "postgres://user:pass@host:5432/db?sslmode=disable", DB: nil}
	s := app.String()
	if !strings.Contains(s, url.QueryEscape("*****")) {
		t.Errorf("String() did not redact DSN: %s", s)
	}
	if !strings.Contains(s, "DB:nil") {
		t.Errorf("String() did not show DB status: %s", s)
	}
}

func TestDBStatus(t *testing.T) {
	if dbStatus(nil) != "nil" {
		t.Errorf("dbStatus(nil) should return 'nil'")
	}
	// Simulate non-nil DB
	if dbStatus((*gorm.DB)(nil)) != "nil" {
		t.Errorf("dbStatus((*gorm.DB)(nil)) should return 'nil'")
	}
}
