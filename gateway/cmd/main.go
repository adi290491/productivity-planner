package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"syscall"
	"time"

	"github.com/adi290491/productivity-planner/gateway/config"
	"github.com/gin-gonic/gin"
)

func init() {

	setupLogging()
}

func setupLogging() {
	profile := os.Getenv("PROFILE")

	var logger *slog.Logger
	if profile == "prod" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo, AddSource: false,
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	slog.SetDefault(logger)
}

func main() {

	appConfig := config.Load()

	slog.Info("Loading configuration",
		"profile", appConfig.Profile,
		"port", appConfig.Port,
	)

	if appConfig.Profile == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.ForceConsoleColor()
	srv := gin.Default()

	RegisterRoutes(srv, appConfig)
	s := &http.Server{
		Addr:         ":" + appConfig.Port,
		ReadTimeout:  appConfig.ReadTimeout,
		WriteTimeout: appConfig.WriteTimeout,
		Handler:      srv,
	}

	go func() {
		slog.Info("Server starting", "port", appConfig.Port, "profile", appConfig.Profile)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	slog.Info("Shutting down...")

	//Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	<-ctx.Done()
	log.Println("Server Exited")
}
