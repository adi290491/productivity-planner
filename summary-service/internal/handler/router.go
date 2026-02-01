package handler

import (
	"net/http"
)

func RegisterEndpoints(mux *http.ServeMux, h *Handler) {

	mux.HandleFunc("GET /health", h.HealthCheck)
	mux.HandleFunc("GET/ready", h.Ready)

	mux.HandleFunc("GET /summary/daily", h.GetDailySummary)
	mux.HandleFunc("GET /summary/weekly", h.GetWeeklySummary)
}
