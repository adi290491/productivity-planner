package config

import (
	"os"
	"time"

	"gorm.io/gorm"
)

type AppConfig struct {
	DSN          string
	Port         string
	DB           *gorm.DB
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() *AppConfig {
	return &AppConfig{ 
		DSN:          os.Getenv("DB_DSN"),
		Port:         os.Getenv("SUMMARY_SERVICE_PORT"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
