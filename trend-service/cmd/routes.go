package main

import "github.com/gin-gonic/gin"

func RegisterEndpoints(r *gin.Engine, h *Handler) {

	r.GET("/trend/daily", h.GetDailyTrend)
	r.GET("/trend/weekly", h.GetWeeklyTrend)
}
