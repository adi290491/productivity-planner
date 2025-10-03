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

func (d *DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}

func dbStatus(db *gorm.DB) string {
	if db == nil {
		return "nil"
	}
	return "initialized"
}

type AppConfig struct {
	Profile   string
	ProjectID string
	Port      string

	// Pub/Sub Topics
	DailyTopic  string
	WeeklyTopic string

	// Pub/Sub Subscriptions
	DailySubscription  string
	WeeklySubscription string

	DSN          string
	DB           *gorm.DB
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func LoadConfig() (*AppConfig, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	profile := getEnvString("PROFILE", "local")

	log.Printf("Loading configurations for %+s\n", profile)

	dbConfig := &DBConfig{
		Host:     getEnvString("DB_HOSTNAME", ""),
		Port:     getEnvString("DB_PORT", ""),
		DBName:   getEnvString("DB_NAME", ""),
		User:     getEnvString("DB_USERNAME", ""),
		Password: getEnvString("DB_PASSWORD", ""),
		SSLMode:  getEnvString("DB_SSLMODE", ""),
	}

	// Validate required fields
	if dbConfig.Host == "" || dbConfig.Port == "" || dbConfig.DBName == "" ||
		dbConfig.User == "" || dbConfig.Password == "" {
		// log.Fatal("Missing required database configuration. Please ensure DB_HOST, DB_PORT, DB_NAME, DB_USERNAME, and DB_PASSWORD are set")
		return nil, fmt.Errorf("missing required database configuration. Please ensure DB_HOST, DB_PORT, DB_NAME, DB_USERNAME, and DB_PASSWORD are set")
	}

	// Set default for optional fields
	if dbConfig.SSLMode == "" {
		dbConfig.SSLMode = "disable"
	}

	appConfig := &AppConfig{
		Profile:   getEnvString("PROFILE", "local"),
		Port:      getEnvString("PORT", "8087"),
		ProjectID: getEnvString("PROJECT_ID", ""),

		DailyTopic:  getEnvString("PUB_SUB_DAILY_TOPIC", ""),
		WeeklyTopic: getEnvString("PUB_SUB_WEEKLY_TOPIC", ""),

		DailySubscription:  getEnvString("PUB_SUB_DAILY_SUBSCRIPTION", ""),
		WeeklySubscription: getEnvString("PUB_SUB_WEEKLY_SUBSCRIPTION", ""),

		DSN: dbConfig.DSN(),
		DB:  nil,
	}

	if appConfig.ProjectID == "" {
		return nil, fmt.Errorf("PROJECT_ID is required")
	}
	if appConfig.DailyTopic == "" {
		return nil, fmt.Errorf("PUB_SUB_DAILY_TOPIC is required")
	}
	if appConfig.DailySubscription == "" {
		return nil, fmt.Errorf("PUB_SUB_DAILY_SUBSCRIPTION is required")
	}

	log.Printf("AppConfig: %+v", *appConfig)
	log.Println("-------------Exiting application config-------------")
	return appConfig, nil
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
