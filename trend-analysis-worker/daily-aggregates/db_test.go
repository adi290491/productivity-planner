package main

import (
	"os"
	"testing"
)


func TestLoadConfig_WithDSN(t *testing.T) {
	const expectedDSN = "host=localhost user=test dbname=testdb sslmode=disable"
	orig := os.Getenv("DB_DSN")
	defer os.Setenv("DB_DSN", orig)

	os.Setenv("DB_DSN", expectedDSN)
	app, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if app.DSN != expectedDSN {
		t.Errorf("expected DSN %s, got %s", expectedDSN, app.DSN)
	}
}

func TestLoadConfig_WithMissingEnv(t *testing.T) {
	os.Unsetenv("DB_DSN")

	app, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for missing DB_DSN, got nil")
	}

	if app != nil {
		t.Errorf("expected nil app on error, got: %v", app)
	}
}
