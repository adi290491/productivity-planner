package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "net/http/pprof"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/adi290491/productivity-planner/session-service/internal/config"
	"github.com/adi290491/productivity-planner/session-service/internal/handler"
	"github.com/adi290491/productivity-planner/session-service/internal/repository/postgres"
	"github.com/adi290491/productivity-planner/session-service/internal/service"
)

func init() {
	setupLogging()
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := postgres.NewDB(cfg.DSN)
	if err != nil {
		slog.Info("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("Database connection successful")

	// Initialize repository layer
	sessionRepo := postgres.NewSessionRepository(db)

	// Initialize service layer
	sessionService := service.NewSessionService(sessionRepo)

	// Initialize HTTP handler
	h := handler.NewHandler(sessionService)

	// Setup router
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, h)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      mux,
	}

	// Start server in goroutine
	go func() {
		slog.Info("Starting HTTP server",
			"port", cfg.Port,
			"read_timeout", cfg.ReadTimeout,
			"write_timeout", cfg.WriteTimeout)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	<-ctx.Done()
	slog.Info("Server stopped gracefully")
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
			Level:     slog.LevelDebug,
			AddSource: true,
		}))
	}

	slog.SetDefault(logger)

	slog.Info("Logging initialized", "profile", profile)
}
