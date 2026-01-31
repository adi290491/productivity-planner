package handler

import (
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) RegisterEndpoints(mux *http.ServeMux) {

	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("GET /ready", h.Ready)

	mux.HandleFunc("POST /users/signup", h.Signup)
	mux.HandleFunc("POST /users/login", h.Login)
	mux.HandleFunc("POST /users/batch", h.GetUsersBatch)

	slog.Info("Routes registered successfully")
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request
		slog.Info("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		// Call next handler
		next.ServeHTTP(w, r)

		// Log completion
		duration := time.Since(start)
		slog.Info("HTTP Response",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered",
					"error", err,
					"method", r.Method,
					"path", r.URL.Path,
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
