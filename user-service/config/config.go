package config

import (
	"os"
	"time"

	"gorm.io/gorm"
)

type AppConfig struct {
	DSN          string
	JWT_SECRET   string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DB           *gorm.DB
}

func Load() *AppConfig {
	return &AppConfig{
		DSN:          os.Getenv("DB_DSN"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		Port:         os.Getenv("USER_SERVICE_PORT"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
