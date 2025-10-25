package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"strings"

	"github.com/adi290491/productivity-planner/trend-service/trend"
	"github.com/adi290491/productivity-planner/trend-service/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *trend.TrendService
}

func (h *Handler) GetDailyTrend(c *gin.Context) {

	log.Println("Inside Get Daily Trend...")
	userId := strings.TrimSpace(c.GetHeader("X-USER-ID"))
	if userId == "" {
		HandleError(c, fmt.Errorf("user id is missing"), http.StatusUnauthorized)
		return
	}

	days := c.DefaultQuery("days", utils.DEFAULT_DAYS)
	log.Println("No of days:", days)
	dailyTrendResponse, err := h.svc.FetchDailyTrend(userId, days)

	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dailyTrendResponse)

}

func (h *Handler) GetWeeklyTrend(c *gin.Context) {
	log.Println("Inside Get Weekly Trend...")
	userId := strings.TrimSpace(c.GetHeader("X-USER-ID"))
	if userId == "" {
		HandleError(c, fmt.Errorf("user id is missing"), http.StatusUnauthorized)
		return
	}

	weeks := c.DefaultQuery("weeks", utils.DEFAULT_WEEKS)

	weeklyTrendResponse, err := h.svc.FetchWeeklyTrend(userId, weeks)

	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, weeklyTrendResponse)

}

func (h *Handler) GetUnviewedTrendsCount(c *gin.Context) {
	log.Println("Inside GetUnviewedTrendsCount...")

	userId := strings.TrimSpace(c.GetHeader("X-USER-ID"))

	if userId == "" {
		HandleError(c, fmt.Errorf("user id is missing"), http.StatusUnauthorized)
		return
	}

	counts, err := h.svc.GetUnviewedTrendsCount(userId)

	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, counts)

}

func (h *Handler) MarkTrendsAsViewed(c *gin.Context) {
	log.Println("Inside MarkTrendsAsViewed...")

	userId := strings.TrimSpace(c.GetHeader("X-USER-ID"))

	if userId == "" {
		HandleError(c, fmt.Errorf("user id is missing"), http.StatusUnauthorized)
		return
	}

	trendType := c.Query("type")
	if trendType != "daily" && trendType != "weekly" {
		HandleError(c, fmt.Errorf("invalid trend type"), http.StatusBadRequest)
		return
	}

	err := h.svc.MarkTrendsAsViewed(userId, trendType)

	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "trend-service",
		"timestamp": time.Now(),
		"profile":   "production",
	})
}

func (h *Handler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
