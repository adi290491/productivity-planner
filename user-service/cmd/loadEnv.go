package main

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Try to resolve env file based on known folder
	// rootPath, err := filepath.Abs(filepath.Join("../", "user-service", ".env"))
	// if err != nil {
	// 	log.Println("Failed to build .env path:", err)
	// 	return
	// }
	rootPath := "/secrets/env-vars"
	err := godotenv.Load(rootPath)
	if err != nil {
		log.Printf("Warning: could not load env from %s\n", rootPath)
	} else {
		log.Printf(".env loaded from %s\n", rootPath)
	}
}
