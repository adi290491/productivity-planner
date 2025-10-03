package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adi290491/productivity-planner/notification-service/config"
	"github.com/adi290491/productivity-planner/notification-service/notification"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	InitDB(config)

	log.Printf("Starting notification service with profile: %s", config.Profile)

	notificationService, err := notification.NewNotificationService(config)
	if err != nil {
		log.Fatalf("Failed to create notification service: %v", err)
	}

	defer notificationService.Close()

	// Create HTTP Server
	server := NewServer(config, notificationService)

	if config.Profile == "production" {
		log.Println("Production mode: Processing daily trends once and exiting")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)

		defer cancel()

		if err := notificationService.ProcessDailyTrendNotifications(ctx); err != nil {
			log.Fatalf("Failed to process daily trend notification: %v", err)
		}

		stats := notificationService.GetStats()
		log.Printf("Processing completed. Stats: %+v", stats)
		return
	}

	log.Println("Development mode: Starting HTTP server for testing")

	// Setup graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutdown signal received")
		notificationService.Close()
		os.Exit(0)
	}()

	log.Printf("Health check: http://localhost:%s/health", config.Port)
	log.Printf("Manual trigger: POST http://localhost:%s/process/daily", config.Port)
	log.Printf("Stats: http://localhost:%s/stat", config.Port)
	log.Printf("Get user notifications: GET http://localhost:%s/notifications?user_id=<uuid>", config.Port)

	if err := server.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
