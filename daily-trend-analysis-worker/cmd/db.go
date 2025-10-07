package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"cloud.google.com/go/pubsub/v2"
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

func dbStatus(db *gorm.DB) string {
	if db == nil {
		return "nil"
	}
	return "initialized"
}

type Application struct {
	DB             *gorm.DB
	DSN            string
	USER_BATCH_URL string
	PROJECT_ID     string
	PUB_SUB_TOPIC  string
	PubSubClient   *pubsub.Client
}

func (c *Application) String() string {
	return "AppConfig{" +
		"DSN:" + redactDSN(c.DSN) +
		", DB:" + dbStatus(c.DB) +
		"}"
}

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

	projectId := os.Getenv("PROJECT_ID")

	// ctx := context.Background()
	// client, err := pubsub.NewClient(ctx, projectId)

	// if err != nil {
	// 	log.Fatalf("Failed to create pubsub client: %v", err)
	// }

	// defer client.Close()
	appConfig := &Application{
		DSN:            dbConfig.DSN(),
		USER_BATCH_URL: os.Getenv("USER_BATCH_URL"),
		PROJECT_ID:     projectId,
		PUB_SUB_TOPIC:  os.Getenv("PUB_SUB_TOPIC"),
		// PubSubClient:   client,
		DB: nil, // DB connection can be set up elsewhere
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
