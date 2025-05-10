package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"productivity-planner/gateway/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}
	// LoadEnv()
}

func main() {

	srv := gin.Default()

	appConfig := config.Load()

	RegisterRoutes(srv)
	log.Println("Port:", appConfig.Port)

	s := &http.Server{
		Addr:         ":" + appConfig.Port,
		ReadTimeout:  appConfig.ReadTimeout,
		WriteTimeout: appConfig.WriteTimeout,
		Handler:      srv,
	}

	// srv.Run(":8000")

	go func() {
		log.Println("Server running on port", appConfig.Port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down...")

	//Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server Exiting")
}
