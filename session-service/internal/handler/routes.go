package handler

import (
	"net/http"
	"net/http/pprof"
	"os"
)

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	// Health checks
	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("GET /ready", h.ReadyCheck)

	// Session endpoints
	mux.HandleFunc("POST /sessions/v1/start-session", h.StartSession)
	mux.HandleFunc("PATCH /sessions/v1/stop-session", h.StopSession)

	if os.Getenv("PROFILE") != "prod" || os.Getenv("PROFILE") != "production" {
		mux.HandleFunc("/debug/pprof", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}
}
