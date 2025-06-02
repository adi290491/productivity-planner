package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"productivity-planner/summary-service/config"
	models "productivity-planner/summary-service/model"

	"productivity-planner/summary-service/summary"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	// LoadEnv()
}

func main() {

	appConfig := config.Load()

	InitDB(appConfig)

	router := gin.Default()

	svc := &summary.SummaryService{
		Repo: &models.PostgresRepository{
			DB: appConfig.DB,
		},
	}
	handler := &Handler{
		Svc: svc,
	}

	RegisterEndpoints(router, handler)

	server := &http.Server{
		Addr:         ":" + appConfig.Port,
		ReadTimeout:  appConfig.ReadTimeout,
		WriteTimeout: appConfig.WriteTimeout,
		Handler:      router,
	}

	go func() {
		log.Println("Server running on port", appConfig.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// graceful shutdown after 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server Exiting")
}
