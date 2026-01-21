package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adi290491/productivity-planner/user-service/config"
	"github.com/adi290491/productivity-planner/user-service/models"
	"github.com/adi290491/productivity-planner/user-service/user"
	"github.com/adi290491/productivity-planner/user-service/utils"
	"github.com/gin-gonic/gin"
)

func init() {
	setupLogging()
}

func main() {
	appConfig, err := config.Load()
	if err != nil {
		slog.Error("Application stopped due to error", "error", err)
		os.Exit(1)
	}

	if appConfig.Profile == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	slog.Info("Starting user-service",
		"profile", appConfig.Profile,
		"port", appConfig.Port,
	)

	// CRITICAL FIX: Initialize database BEFORE starting server
	slog.Info("Starting database initialization...")
	if err := InitDB(appConfig); err != nil {
		slog.Error("Database initialization failed", "error", err)
		os.Exit(1)
	}
	dbInitialized.Store(true)
	slog.Info("Database initialization complete")

	// Create service with initialized DB
	svc := &user.UserService{
		Repo: &models.PostgresRepository{
			DB: appConfig.DB,
		},
	}

	handler := &Handler{
		Svc: svc,
		JwtUtil: utils.JWTUtil{
			Secret: []byte(appConfig.JWT_SECRET),
		},
	}

	// Create server and register ALL endpoints BEFORE starting
	server := gin.Default()
	RegisterEndpoints(server, handler)

	s := &http.Server{
		Addr:         ":" + appConfig.Port,
		ReadTimeout:  appConfig.ReadTimeout,
		WriteTimeout: appConfig.WriteTimeout,
		Handler:      server,
	}

	go func() {
		slog.Info("Server listening", "port", appConfig.Port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped gracefully")
}

func setupLogging() {
	profile := os.Getenv("PROFILE")

	var logger *slog.Logger
	if profile == "prod" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	slog.SetDefault(logger)
}
