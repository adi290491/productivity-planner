package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
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

type Application struct {
	DB  *gorm.DB
	DSN string
}

func LoadConfig() (*Application, error) {

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
		// log.Fatal("Missing required database configuration. Please ensure DB_HOST, DB_PORT, DB_NAME, DB_USERNAME, and DB_PASSWORD are set")
		return nil, fmt.Errorf("missing required database configuration. Please ensure DB_HOST, DB_PORT, DB_NAME, DB_USERNAME, and DB_PASSWORD are set")
	}

	// Set default for optional fields
	if dbConfig.SSLMode == "" {
		dbConfig.SSLMode = "disable"
	}
	appConfig := &Application{
		DSN: dbConfig.DSN(),
		DB:  nil, // DB connection can be set up elsewhere
	}

	log.Println("-------------Exiting application config-------------")
	return appConfig, nil
}

func (a *Application) InitDB() error {

	db, err := gorm.Open(postgres.Open(a.DSN), &gorm.Config{})

	if err != nil {
		log.Fatalf("init db error: %v", err)
	}

	sqldb, err := db.DB()
	if err != nil {
		return err
	}

	if err = sqldb.Ping(); err != nil {
		return fmt.Errorf("could not connect to db: %v", err)
	}

	a.DB = db
	return nil
}
