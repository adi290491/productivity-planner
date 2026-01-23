package config

import (
	"fmt"
	"log/slog"
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
	slog.Info("Loading configurations", "profile", profile)

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
		slog.Error("Missing required database configuration",
			"required", []string{"DB_HOSTNAME", "DB_PORT", "DB_NAME", "DB_USERNAME", "DB_PASSWORD"})
		os.Exit(1)
	}

	// Set default for optional fields
	if dbConfig.SSLMode == "" {
		dbConfig.SSLMode = "disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		slog.Warn("PORT not specified, using default", "port", port)
	}

	cfg := &Config{
		DSN:          dbConfig.DSN(),
		Port:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("Configuration loaded",
		"port", cfg.Port,
		"db_host", dbConfig.Host,
		"db_port", dbConfig.Port,
		"db_name", dbConfig.DBName, // ← Check this value
		"dsn", cfg.DSN) // ← See the full DSN
	return cfg
}

func (d *DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}
