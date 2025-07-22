package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"time"

	"github.com/adi290491/productivity-planner/task-service/config"
	"github.com/adi290491/productivity-planner/task-service/models"
	"github.com/adi290491/productivity-planner/task-service/session"
	"github.com/gin-gonic/gin"
)

func init() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("No .env file found")
	// }
	// LoadEnv()
}

func main() {

	appConfig := config.Load()

	InitDB(appConfig)

	server := gin.Default()

	svc := &session.SessionService{
		Repo: &models.PostgresRepository{
			DB: appConfig.DB,
		},
	}

	handler := Handler{Svc: svc}

	RegisterEndpoints(server, &handler)

	s := &http.Server{
		Addr:         ":" + appConfig.Port,
		ReadTimeout:  appConfig.ReadTimeout,
		WriteTimeout: appConfig.WriteTimeout,
		Handler:      server,
	}

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
