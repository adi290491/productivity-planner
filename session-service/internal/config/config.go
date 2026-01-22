package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
}

type Config struct {
	DSN          string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() *Config {
	// Load .env file (ignore error if not present)
	_ = godotenv.Load()

	profile := os.Getenv("PROFILE")
	log.Printf("Loading configurations for %s\n", profile)

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
		log.Fatal("Missing required database configuration. Please ensure DB_HOSTNAME, DB_PORT, DB_NAME, DB_USERNAME, and DB_PASSWORD are set")
	}

	// Set default for optional fields
	if dbConfig.SSLMode == "" {
		dbConfig.SSLMode = "disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("PORT not specified, using default: %s", port)
	}

	cfg := &Config{
		DSN:          dbConfig.DSN(),
		Port:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Configuration loaded successfully")
	return cfg
}

func (d *DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}
