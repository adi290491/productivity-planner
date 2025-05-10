package main

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Try to resolve env file based on known folder
	rootPath, err := filepath.Abs(filepath.Join("../", "session-service", ".env"))
	if err != nil {
		log.Println("Failed to build .env path:", err)
		return
	}

	err = godotenv.Load(rootPath)
	if err != nil {
		log.Printf("Warning: could not load env from %s\n", rootPath)
	} else {
		log.Printf(".env loaded from %s\n", rootPath)
	}
}
