package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_ValidEnvVars(t *testing.T) {
	os.Setenv("PORT", "9000")
	os.Setenv("JWT_SECRET", "testsecret")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg := Load()
	if cfg.Port != "9000" {
		t.Errorf("expected port 9000, got %s", cfg.Port)
	}
	if cfg.JWT_SECRET != "testsecret" {
		t.Errorf("expected JWT_SECRET testsecret, got %s", cfg.JWT_SECRET)
	}
	if cfg.ReadTimeout != 10*time.Second {
		t.Errorf("expected ReadTimeout 10s, got %v", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 10*time.Second {
		t.Errorf("expected WriteTimeout 10s, got %v", cfg.WriteTimeout)
	}
}

func TestLoad_DefaultPort(t *testing.T) {
	os.Unsetenv("PORT")
	os.Setenv("JWT_SECRET", "testsecret")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg := Load()
	if cfg.Port != "8081" {
		t.Errorf("expected default port 8081, got %s", cfg.Port)
	}
}

func TestAppConfig_String(t *testing.T) {
	cfg := &AppConfig{
		Port:         "9000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		JWT_SECRET:   "testsecret",
	}
	str := cfg.String()
	if str == "" || str == "AppConfig{}" {
		t.Errorf("unexpected string output: %s", str)
	}
}
