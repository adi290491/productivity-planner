package main

import "github.com/gin-gonic/gin"

func RegisterEndpoints(r *gin.Engine, h *Handler) {

	r.GET("/health", h.HealthCheck)
	r.GET("/ready", h.Ready)

	r.GET("/trend/daily", h.GetDailyTrend)
	r.GET("/trend/weekly", h.GetWeeklyTrend)
	r.GET("/trend/unviewed", h.GetUnviewedTrendsCount)
	r.POST("/trend/mark-viewed", h.MarkTrendsAsViewed)
}
