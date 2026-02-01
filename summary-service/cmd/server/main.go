package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "net/http/pprof"

	"github.com/adi290491/productivity-planner/summary-service/internal/config"
	"github.com/adi290491/productivity-planner/summary-service/internal/handler"
	"github.com/adi290491/productivity-planner/summary-service/internal/repository"
	summary "github.com/adi290491/productivity-planner/summary-service/internal/service"
)

func init() {
	setupLogging()
}

func main() {

	appConfig, err := config.Load()

	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting summary-service",
		"port", appConfig.Port,
		"profile", appConfig.Profile,
		"readTimeout", appConfig.ReadTimeout,
		"writeTimeout", appConfig.WriteTimeout,
	)

	repo, err := repository.NewPostgresRepository(appConfig.DSN)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer repo.Close()

	svc := summary.NewSummaryService(repo)

	h := handler.NewHandler(svc)

	mux := http.NewServeMux()
	handler.RegisterEndpoints(mux, h)

	server := &http.Server{
		Addr:         ":" + appConfig.Port,
		ReadTimeout:  appConfig.ReadTimeout,
		WriteTimeout: appConfig.WriteTimeout,
		Handler:      mux,
	}

	go func() {
		slog.Info("Server listening", "port", appConfig.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited successfully")
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
