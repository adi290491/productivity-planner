package main

import (
	"testing"

	"github.com/adi290491/productivity-planner/user-service/config"
)

func TestInitDB_NilConfig(t *testing.T) {
	err := InitDB(nil)
	if err == nil {
		t.Error("Expected error for nil config, got nil")
	}
	if err.Error() != "InitDB: appConfig is nil" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestInitDB_EmptyDSN(t *testing.T) {
	appConfig := &config.AppConfig{
		DSN: "",
	}

	err := InitDB(appConfig)
	if err == nil {
		t.Error("Expected error for empty DSN, got nil")
	}
	if err.Error() != "InitDB: DSN is empty" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}
