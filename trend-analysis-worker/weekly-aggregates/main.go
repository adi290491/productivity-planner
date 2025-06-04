package main

import (
	"log"
	"os"

	"gorm.io/gorm"
)

type Application struct {
	DB  *gorm.DB
	DSN string
}

func main() {
	app := LoadConfig()
	log.SetOutput(os.Stdout)
	log.Println("Weekly trend job started")

	app.InitDB()
	app.FetchWeeklyTrend()
}
