package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Try to resolve env file based on known folder
	// rootPath, err := filepath.Abs(filepath.Join("../", "user-service", ".env"))
	// if err != nil {
	// 	log.Println("Failed to build .env path:", err)
	// 	return
	// }
	// rootPath := "/secrets/env-vars"
	// err := godotenv.Load(rootPath)
	// if err != nil {
	// 	log.Printf("Warning: could not load env from %s\n", rootPath)
	// } else {
	// 	log.Printf(".env loaded from %s\n", rootPath)
	// }
	secretPath := "./secrets/env-var"

	// Local .env fallback path
	localPath := ".env"

	// Attempt to load from secret mount first
	if _, err := os.Stat(secretPath); err == nil {
		err := godotenv.Load(secretPath)
		if err != nil {
			log.Printf("Failed to load env from secret mount %s: %v", secretPath, err)
		} else {
			log.Printf("Loaded env from secret mount %s", secretPath)
			return
		}
	}

	// Fallback to local .env file for dev/test
	if _, err := os.Stat(localPath); err == nil {
		err := godotenv.Load(localPath)
		if err != nil {
			log.Printf("Failed to load local .env: %v", err)
		} else {
			log.Println("Loaded local .env")
		}
	} else {
		log.Println(" No .env file found locally or in secrets. Proceeding with existing env vars.")
	}
}
