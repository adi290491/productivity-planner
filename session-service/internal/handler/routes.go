package handler

import (
	"net/http"
)

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	// Health checks
	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("GET /ready", h.ReadyCheck)

	// Session endpoints
	mux.HandleFunc("POST /sessions/v1/start-session", h.StartSession)
	mux.HandleFunc("PATCH /sessions/v1/stop-session", h.StopSession)
}
