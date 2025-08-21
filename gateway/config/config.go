package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	JWT_SECRET   string
}

func (c *AppConfig) String() string {
	return fmt.Sprintf("AppConfig{Port:%s, ReadTimeout:%s, WriteTimeout:%s}", c.Port, c.ReadTimeout, c.WriteTimeout)
}

func Load() *AppConfig {

	_ = godotenv.Load()

	profile := os.Getenv("PROFILE")
	log.Printf("Loading configurations for %+s\n", profile)

	appConfig := &AppConfig{
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		Port:         os.Getenv("PORT"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Validate required application settings
	if appConfig.JWT_SECRET == "" {
		log.Fatal("Missing required JWT_SECRET environment variable")
	}

	if appConfig.Port == "" {
		appConfig.Port = "8081" // Set default port
		log.Printf("PORT not specified, using default: %s", appConfig.Port)
	}

	log.Println("-------------Finished application config-------------")
	return appConfig

}
