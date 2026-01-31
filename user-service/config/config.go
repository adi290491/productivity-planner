package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Profile      string
	Port         string
	Database     DatabaseConfig
	JWT          JWTConfig
	Server       ServerConfig
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
	DSN      string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret []byte
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	profile := getEnv("PROFILE", "local")

	dbConfig := DatabaseConfig{
		Host:     getEnv("DB_HOSTNAME", ""),
		Port:     getEnv("DB_PORT", ""),
		Name:     getEnv("DB_NAME", ""),
		User:     getEnv("DB_USERNAME", ""),
		Password: getEnv("DB_PASSWORD", ""),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Validate required database fields
	if dbConfig.Host == "" || dbConfig.Port == "" || dbConfig.Name == "" ||
		dbConfig.User == "" || dbConfig.Password == "" {
		return nil, fmt.Errorf("missing required database configuration. Please ensure DB_HOSTNAME, DB_PORT, DB_NAME, DB_USERNAME, and DB_PASSWORD are set")
	}

	// Build DSN
	dbConfig.DSN = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.SSLMode,
	)

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, fmt.Errorf("missing required JWT_SECRET environment variable")
	}

	port := getEnv("PORT", "8081")

	config := &Config{
		Profile:  profile,
		Port:     port,
		Database: dbConfig,
		JWT: JWTConfig{
			Secret: []byte(jwtSecret),
		},
		Server: ServerConfig{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	return config, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
