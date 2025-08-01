package config

import (
	"log"
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

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set by Cloud Run.")
	}

	return &AppConfig{
		DSN:          os.Getenv("DB_DSN"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		Port:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
