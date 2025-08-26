package main

import (
	"os"
	"strings"
	"testing"
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
	if !strings.Contains(dsn, "testuser:testpass@localhost:5432/testdb?sslmode=disable") {
		t.Errorf("DSN not constructed properly: %s", dsn)
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
