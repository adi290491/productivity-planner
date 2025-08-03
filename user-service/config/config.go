package config

import (
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

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
		", JWT_SECRET:" + c.JWT_SECRET +
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

	content := readFromFile("/secrets/env-var")

	log.Println("Content:", content)

	port := os.Getenv("PORT")
	log.Println("PORT:", port)

	return &AppConfig{
		DSN:          os.Getenv("DB_DSN"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		Port:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func readFromFile(path string) string {
	data, err := os.ReadFile(path)

	if err != nil {
		log.Printf("Warning: could not read file at %s\n", err)
		return ""
	}

	return strings.TrimSpace(string(data))

}
