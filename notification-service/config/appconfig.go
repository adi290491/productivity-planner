package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

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
}

func LoadConfig() (*AppConfig, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	appConfig := &AppConfig{
		Profile:   getEnvString("PROFILE", "local"),
		Port:      getEnvString("PORT", "8087"),
		ProjectID: getEnvString("PROJECT_ID", ""),

		DailyTopic:  getEnvString("PUB_SUB_DAILY_TOPIC", ""),
		WeeklyTopic: getEnvString("PUB_SUB_WEEKLY_TOPIC", ""),

		DailySubscription:  getEnvString("PUB_SUB_DAILY_SUBSCRIPTION", ""),
		WeeklySubscription: getEnvString("PUB_SUB_WEEKLY_SUBSCRIPTION", ""),
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
