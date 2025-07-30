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
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("TREND_SERVICE_PORT")
	}
	return &AppConfig{
		DSN:          os.Getenv("DB_DSN"),
		Port:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
