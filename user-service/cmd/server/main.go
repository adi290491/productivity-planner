package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/adi290491/productivity-planner/user-service/config"
	"github.com/adi290491/productivity-planner/user-service/internal/handler"
	"github.com/adi290491/productivity-planner/user-service/internal/repository"
	"github.com/adi290491/productivity-planner/user-service/internal/service"
	"github.com/adi290491/productivity-planner/user-service/pkg/util"
)

func init() {
	setupLogging()
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting user-service", "profile", cfg.Profile, "port", cfg.Port)

	db, err := initDB(cfg)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Database initialized successfully")

	// Create service with initialized DB
	userRepo := repository.NewPostgresRepository(db)
	userService := service.NewUserService(userRepo)
	jwtUtil := util.NewJWTUtil(cfg.JWT.Secret)
	handler := handler.NewHandler(userService, jwtUtil)

	mux := http.NewServeMux()
	handler.RegisterEndpoints(mux)

	s := &http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		Handler:      mux,
	}

	go func() {
		slog.Info("Server listening", "port", cfg.Port)
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

// NewDB creates a new database connection
func initDB(cfg *config.Config) (*sql.DB, error) {
	slog.Info("Initializing database connection")

	db, err := sql.Open("pgx", cfg.Database.DSN)

	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
