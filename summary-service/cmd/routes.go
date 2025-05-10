package main

import "github.com/gin-gonic/gin"

func RegisterEndpoints(r *gin.Engine, h *Handler) {

	r.GET("/summary/daily", h.GetDailySummary)
	r.GET("/summary/weekly", h.GetWeeklySummary)
}
