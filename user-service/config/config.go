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
	DbName   string
	User     string
	Password string
	SSLMode  string
}

type AppConfig struct {
	DSN          string
	JWT_SECRET   string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DB           *gorm.DB
}

// String implements the fmt.Stringer interface for AppConfig.
func (c *AppConfig) String() string {
	return "AppConfig{" +
		"DSN:" + c.DSN +
		", Port:" + c.Port +
		", ReadTimeout:" + c.ReadTimeout.String() +
		", WriteTimeout:" + c.WriteTimeout.String() +
		", DB:" + dbStatus(c.DB) +
		"}"
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
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DbName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Validate required fields
	if dbConfig.Host == "" || dbConfig.Port == "" || dbConfig.DbName == "" ||
		dbConfig.User == "" || dbConfig.Password == "" {
		log.Fatal("Missing required database configuration. Please ensure DB_HOST, DB_PORT, DB_NAME, DB_USERNAME, and DB_PASSWORD are set")
	}

	// Set default for optional fields
	if dbConfig.SSLMode == "" {
		dbConfig.SSLMode = "disable"
	}
	appConfig := &AppConfig{
		DSN:          dbConfig.DSN(),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		Port:         os.Getenv("PORT"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		DB:           nil, // DB connection can be set up elsewhere
	}

	// Validate required application settings
	if appConfig.JWT_SECRET == "" {
		log.Fatal("Missing required JWT_SECRET environment variable")
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
		d.User, d.Password, d.Host, d.Port, d.DbName, d.SSLMode,
	)
}
