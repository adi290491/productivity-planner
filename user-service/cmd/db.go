package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/adi290491/productivity-planner/user-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(appConfig *config.AppConfig) error {
	slog.Info("Initializing Database")

	if appConfig == nil {
		return fmt.Errorf("InitDB: appConfig is nil")
	}

	if appConfig.DSN == "" {
		return fmt.Errorf("InitDB: DSN is empty")
	}

	db, err := gorm.Open(postgres.Open(appConfig.DSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}

	// Optimize connection pool for Cloud Run cold starts
	sqlDB.SetMaxIdleConns(2)                   // Reduce idle connections
	sqlDB.SetMaxOpenConns(10)                  // Limit max connections
	sqlDB.SetConnMaxLifetime(time.Hour)        // Reuse connections for 1 hour
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Close idle after 10min

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	slog.Info("Database connection successful")
	appConfig.DB = db
	return nil
}
