package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adi290491/productivity-planner/gateway/config"
	"github.com/adi290491/productivity-planner/gateway/ratelimiter"
	"github.com/adi290491/productivity-planner/gateway/router"
)

func init() {
	setupLogging()
}

func main() {

	slog.Info("Environment variables",
		"PROFILE", os.Getenv("PROFILE"),
		"PORT", os.Getenv("PORT"),
		"USER_SERVICE_URL", os.Getenv("USER_SERVICE_URL"),
		"FRONTEND_ORIGIN", os.Getenv("FRONTEND_ORIGIN"),
		"FRONTEND_ORIGIN_2", os.Getenv("FRONTEND_ORIGIN_2"),
	)

	appConfig := config.Load()

	slog.Info("Starting gateway",
		"profile", appConfig.Profile,
		"port", appConfig.Port,
	)

	rateLimitConfig := ratelimiter.LoadConfig()
	if err := rateLimitConfig.Validate(); err != nil {
		slog.Error("Invalid rate limited configureation", "error", err)
		os.Exit(1)
	}

	slog.Info("Rate limited configured", "config", rateLimitConfig.String())

	// Create rate limiter manager
	var rateLimiterMgr *ratelimiter.RateLimiterManager
	if rateLimitConfig.Enabled {
		rateLimiterMgr = ratelimiter.NewRateLimiterManager(
			rateLimitConfig.Capacity,
			rateLimitConfig.RefillRate,
			rateLimitConfig.CleanupInterval,
			rateLimitConfig.TTL,
		)
		defer rateLimiterMgr.Shutdown()
		slog.Info("Rate limiter enabled")
	} else {
		slog.Warn("Rate limiter disabled")

		rateLimiterMgr = ratelimiter.NewRateLimiterManager(1, 1, 1*time.Minute, 1*time.Minute)
		defer rateLimiterMgr.Shutdown()
	}

	handler := router.NewRouter(appConfig, rateLimiterMgr, rateLimitConfig)

	// Cretae HTTP server
	server := &http.Server{
		Addr:         ":" + appConfig.Port,
		Handler:      handler,
		ReadTimeout:  appConfig.ReadTimeout,
		WriteTimeout: appConfig.WriteTimeout,
	}

	// Start server
	go func() {
		slog.Info("Server listening",
			"port", appConfig.Port,
			"address", "http://localhost:"+appConfig.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	slog.Info("Shutting down server...")

	//Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
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
			Level: slog.LevelInfo, AddSource: false,
		}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	slog.SetDefault(logger)
}
