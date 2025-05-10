package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"productivity-planner/task-service/config"
	"productivity-planner/task-service/controller"
	"productivity-planner/task-service/models"
	"productivity-planner/task-service/repository"
	"productivity-planner/task-service/session"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("No .env file found")
	// }
	LoadEnv()
}

func main() {

	appConfig := config.Load()

	repository.InitDB(appConfig)

	server := gin.Default()

	svc := session.NewSessionService(models.NewPostgresRepository(appConfig.DB))
	handler := controller.NewHandler(svc)

	controller.RegisterEndpoints(server, handler)

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
