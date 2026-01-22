package handler

import (
	"net/http"
)

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	// Health checks
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/ready", h.ReadyCheck)

	// Session endpoints
	mux.HandleFunc("/sessions/v1/start-session", h.StartSession)
	mux.HandleFunc("/sessions/v1/stop-session", h.StopSession)
}
