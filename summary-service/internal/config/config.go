package config

import (
	"fmt"
	"log/slog"
	"net/url"
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

func (d *DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}

type AppConfig struct {
	DSN          string
	Port         string
	Profile      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c *AppConfig) String() string {
	return fmt.Sprintf("Config{DSN:%s, Port:%s, ReadTimeout:%s, WriteTimeout:%s, Profile:%s}",
		redactDSN(c.DSN), c.Port, c.ReadTimeout, c.WriteTimeout, c.Profile)
}

// redactDSN masks the password in a DSN string. Returns original DSN on parse error or if no user info.
func redactDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil || u == nil || u.User == nil {
		return dsn
	}
	pw, _ := u.User.Password()
	if pw == "" {
		return dsn
	}
	u.User = url.UserPassword(u.User.Username(), "*****")
	return u.String()
}

func Load() (*AppConfig, error) {

	// Load .env file if it exists, but don't fail if it doesn't
	if err := godotenv.Load(); err != nil {
		slog.Debug("No .env file found, using system environment variables")
	}

	profile := os.Getenv("PROFILE")
	if profile == "" {
		profile = "local"
	}

	slog.Info("Loading configurations", "profilr", profile)

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
		// log.Fatal("Missing required database configuration. Please ensure DB_HOST, DB_PORT, DB_NAME, DB_USERNAME, and DB_PASSWORD are set")
		return nil, fmt.Errorf("missing required database configuration. Please ensure DB_HOST, DB_PORT, DB_NAME, DB_USERNAME, and DB_PASSWORD are set")
	}

	// Set default for optional fields
	if dbConfig.SSLMode == "" {
		dbConfig.SSLMode = "disable"
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8082"
		slog.Debug("PORT not specified, using default", "port", port)
	}

	appConfig := &AppConfig{
		DSN:          dbConfig.DSN(),
		Port:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Profile:      profile,
	}

	slog.Info("Configuration loaded successfully",
		"port", appConfig.Port,
		"profile", appConfig.Profile,
		"db_host", dbConfig.Host,
		"db_name", dbConfig.DBName,
	)
	return appConfig, nil
}
