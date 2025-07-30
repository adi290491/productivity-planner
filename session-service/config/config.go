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

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("SESSION_SERVICE_PORT")
	}
	return &AppConfig{
		DSN:          os.Getenv("DB_DSN"),
		Port:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
