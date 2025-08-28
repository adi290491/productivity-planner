package main

import "github.com/gin-gonic/gin"

func RegisterEndpoints(r *gin.Engine, h *Handler) {

	r.GET("/trend/daily", h.GetDailyTrend)
	r.GET("/trend/weekly", h.GetWeeklyTrend)
}

/*
Lets connect with the backend. Create a trend-analysis.tsx inside api folder
Endpoint for dialy analysis:
/trend/daily

Endpoint for weekly analysis:
/trend/weekly
*/