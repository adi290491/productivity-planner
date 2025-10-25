package main

import "github.com/gin-gonic/gin"

func RegisterEndpoints(r *gin.Engine, h *Handler) {

	r.GET("/health", h.HealthCheck)
	r.GET("/ready", h.Ready)

	r.GET("/summary/daily", h.GetDailySummary)
	r.GET("/summary/weekly", h.GetWeeklySummary)
}
