package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	JWT_SECRET   string
}

// func Load() *AppConfig {

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = os.Getenv("GATEWAY_PORT")
// 	}

// 	return &AppConfig{
// 		Port:         port,
// 		ReadTimeout:  10 * time.Second,
// 		WriteTimeout: 10 * time.Second,
// 	}
// }

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
		Profile   string
		Port      string
		JWTSecret string `yaml:"jwt_secret"`
	}

	var rc RawConfig

	err = yaml.Unmarshal(configData, &rc)

	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	rc.Port = os.ExpandEnv(rc.Port)
	if profile == "prod" {
		log.Println("Fetching env variables for prod....")
		rc.Port = os.ExpandEnv(rc.Port)
		rc.JWTSecret = os.ExpandEnv(rc.JWTSecret)
	}

	appConfig := &AppConfig{
		JWT_SECRET:   rc.JWTSecret,
		Port:         rc.Port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
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
