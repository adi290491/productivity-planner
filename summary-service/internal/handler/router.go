package handler

import (
	"net/http"
	"net/http/pprof"
	"os"
)

func RegisterEndpoints(mux *http.ServeMux, h *Handler) {

	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("GET /ready", h.Ready)

	mux.HandleFunc("GET /summary/daily", h.GetDailySummary)
	mux.HandleFunc("GET /summary/weekly", h.GetWeeklySummary)

	if os.Getenv("PROFILE") != "prod" || os.Getenv("PROFILE") != "production" {
		mux.HandleFunc("/debug/pprof", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}
}
