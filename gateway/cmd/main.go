package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"productivity-planner/gateway/config"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {

	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("No .env file found")
	// }
	// LoadEnv()
}

func main() {
	gin.ForceConsoleColor()
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
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	log.Println("Shutting down...")

	//Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	<-ctx.Done()
	log.Println("timeout of 5 seconds")
	log.Println("Server Exiting")
}
