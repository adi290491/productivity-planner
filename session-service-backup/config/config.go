package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
}

type AppConfig struct {
	DSN          string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DB           *gorm.DB
}

func dbStatus(db *gorm.DB) string {
	if db == nil {
		return "nil"
	}
	return "initialized"
}

func Load() *AppConfig {

	_ = godotenv.Load()

	profile := os.Getenv("PROFILE")

	log.Printf("Loading configurations for %+s\n", profile)

	dbConfig := &DBConfig{
		Host:     os.Getenv("DB_HOSTNAME"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Validate required fields
	if dbConfig.Host == "" || dbConfig.Port == "" || dbConfig.DBName == "" ||
		dbConfig.User == "" || dbConfig.Password == "" {
		log.Fatal("Missing required database configuration. Please ensure DB_HOST, DB_PORT, DB_NAME, DB_USERNAME, and DB_PASSWORD are set")
	}

	// Set default for optional fields
	if dbConfig.SSLMode == "" {
		dbConfig.SSLMode = "disable"
	}
	appConfig := &AppConfig{
		DSN:          dbConfig.DSN(),
		Port:         os.Getenv("PORT"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		DB:           nil, // DB connection can be set up elsewhere
	}

	if appConfig.Port == "" {
		appConfig.Port = "8080" // Set default port
		log.Printf("PORT not specified, using default: %s", appConfig.Port)
	}

	log.Println("-------------Exiting application config-------------")
	return appConfig

}

func (d *DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}
