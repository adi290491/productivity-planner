package main

import "github.com/gin-gonic/gin"

func RegisterEndpoints(r *gin.Engine, h *Handler) {

	r.GET("/trend/daily", h.GetDailyTrend)
	r.GET("/trend/weekly", h.GetWeeklyTrend)
	r.GET("/trend/unViewed", h.GetUnviewedTrendsCount)
	r.POST("/trend/mark-viewed", h.MarkTrendsAsViewed)
}

/*
Lets connect with the backend. Create a trend-analysis.tsx inside api folder
Endpoint for dialy analysis:
/trend/daily

Endpoint for weekly analysis:
/trend/weekly
*/
