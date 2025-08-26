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

	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	profile := os.Getenv("PROFILE")

	log.Printf("Loading configurations for %+s\n", profile)
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("DB_DSN environment variable is required")
	}
	return &Application{
		DSN: dsn,
	}, nil
}

func (a *Application) InitDB() {

	db, err := gorm.Open(postgres.Open(a.DSN), &gorm.Config{})

	if err != nil {
		log.Fatalf("init db error: %v", err)
	}

	sqldb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err = sqldb.Ping(); err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	a.DB = db
}
