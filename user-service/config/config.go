package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"sslmode"`
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

	profile := os.Getenv("PROFILE")
	if profile == "" {
		profile = "local"
	}

	configData, err := os.ReadFile(getConfigPath(profile))
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	type RawConfig struct {
		Profile string
		Port    string
		DB      struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Database string `yaml:"database"`
			SSLMode  string `yaml:"sslmode"`
		} `yaml:"db"`
		JWTSecret string `yaml:"jwt_secret"`
	}

	var rc RawConfig

	err = yaml.Unmarshal(configData, &rc)

	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	rc.Port = os.ExpandEnv(rc.Port)
	if profile == "prod" {
		// rc.Port = os.ExpandEnv(rc.Port)
		rc.DB.Username = os.ExpandEnv(rc.DB.Username)
		rc.DB.Password = os.ExpandEnv(rc.DB.Password)
		rc.DB.Host = os.ExpandEnv(rc.DB.Host)
		rc.DB.Port = os.ExpandEnv(rc.DB.Port)
		rc.DB.Database = os.ExpandEnv(rc.DB.Database)
		rc.DB.SSLMode = os.ExpandEnv(rc.DB.SSLMode)
		rc.JWTSecret = os.ExpandEnv(rc.JWTSecret)
	}

	dbConfig := &DBConfig{
		Host:     rc.DB.Host,
		Port:     rc.DB.Port,
		Database: rc.DB.Database,
		User:     rc.DB.Username,
		Password: rc.DB.Password,
		SSLMode:  rc.DB.SSLMode,
	}

	appConfig := &AppConfig{
		DSN:          dbConfig.DSN(),
		JWT_SECRET:   rc.JWTSecret,
		Port:         rc.Port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		DB:           nil, // DB connection can be set up elsewhere
	}

	log.Println("-------------Read application config-------------")
	return appConfig

}

func getConfigPath(profile string) string {
	var configPath string

	dir, _ := os.Getwd()
	log.Printf("Profile : %s\n, CWD: %s", profile, dir)
	// if profile == "prod" {
	// 	configPath = "/config/config.prod.yaml"
	// } else {
	// 	configPath = "/config/config.local.yaml"
	// }

	configPath = fmt.Sprintf("config/config.%s.yaml", profile)

	return configPath
}

func readFromFile(path string) string {
	data, err := os.ReadFile(path)

	if err != nil {
		log.Printf("Warning: could not read file at %s\n", err)
		return ""
	}

	return strings.TrimSpace(string(data))

}

func (d *DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Database, d.SSLMode,
	)
}

/*
 */
