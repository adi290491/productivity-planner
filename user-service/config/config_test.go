package config

import (
	"os"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestDBConfig_DSN(t *testing.T) {

	dbConf := &DBConfig{
		Host:     "localhost",
		Port:     "5432",
		DbName:   "testdb",
		User:     "testuser",
		Password: "testpass",
		SSLMode:  "disable",
	}
	expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	if dbConf.DSN() != expected {
		t.Errorf("expected %s, got %s", expected, dbConf.DSN())
	}
}

func TestAppConfig_String(t *testing.T) {

	appConf := &AppConfig{
		DSN:          "dsn",
		JWT_SECRET:   "secret",
		Port:         "1234",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		DB:           nil,
	}
	str := appConf.String()
	if str == "" || str == "AppConfig{}" {
		t.Errorf("unexpected string output: %s", str)
	}
}

func TestDbStatus(t *testing.T) {

	if dbStatus(nil) != "nil" {
		t.Errorf("expected nil for nil DB")
	}
	// Simulate a non-nil DB (not a real connection)
	db := &gorm.DB{}
	if dbStatus(db) != "initialized" {
		t.Errorf("expected initialized for non-nil DB")
	}
}

func TestLoad_WithEnvVars(t *testing.T) {

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USERNAME", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("JWT_SECRET", "secret")
	os.Setenv("PORT", "1234")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_USERNAME")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("PORT")
	}()

	// SSLMode should default to "disable" if not set
	os.Unsetenv("DB_SSLMODE")

	cfg := Load()
	if cfg.Port != "1234" {
		t.Errorf("expected port 1234, got %s", cfg.Port)
	}
	if cfg.JWT_SECRET != "secret" {
		t.Errorf("expected JWT_SECRET secret, got %s", cfg.JWT_SECRET)
	}
	if cfg.DSN != "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable" {
		t.Errorf("unexpected DSN: %s", cfg.DSN)
	}
	if cfg.DB != nil {
		t.Errorf("expected DB to be nil")
	}
}

// func TestLoad_MissingEnvVars(t *testing.T) {

// 	os.Unsetenv("DB_HOST")
// 	os.Unsetenv("DB_PORT")
// 	os.Unsetenv("DB_NAME")
// 	os.Unsetenv("DB_USERNAME")
// 	os.Unsetenv("DB_PASSWORD")
// 	os.Unsetenv("DB_SSLMODE")
// 	os.Unsetenv("JWT_SECRET")
// 	os.Unsetenv("PORT")

// 	// Expect Load to call log.Fatal due to missing required DB config
// 	// Use a recoverable test for log.Fatal
// 	defer func() {
// 		if r := recover(); r == nil {
// 			t.Errorf("expected panic (log.Fatal) due to missing required env vars, but did not panic")
// 		}
// 	}()
// 	Load()
// }
