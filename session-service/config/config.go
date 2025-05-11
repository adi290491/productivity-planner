package config

import (
	"os"
	"time"

	"gorm.io/gorm"
)

type AppConfig struct {
	DSN          string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DB           *gorm.DB
}

func Load() *AppConfig {
	return &AppConfig{
		DSN:          os.Getenv("DB_DSN"),
		Port:         os.Getenv("SESSION_SERVICE_PORT"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
