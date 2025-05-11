package config

import (
	"os"
	"time"
)

type AppConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() *AppConfig {
	return &AppConfig{
		Port:         os.Getenv("GATEWAY_PORT"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
