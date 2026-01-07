package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"strings"

	"github.com/adi290491/productivity-planner/summary-service/summary"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Svc summary.SummaryServiceInterface
}

func (h *Handler) GetDailySummary(c *gin.Context) {

	userId := strings.TrimSpace(c.GetHeader("X-USER-ID"))
	if userId == "" {
		HandleError(c, fmt.Errorf("user id is missing"), http.StatusUnauthorized)
		return
	}

	queryDate := c.Query("date")
	log.Println("Query Date:", queryDate)
	summaryResponse, err := h.Svc.GetDailySessionSummary(userId, queryDate)

	if err != nil && strings.Contains(err.Error(), "invalid date format") {
		HandleError(c, err, http.StatusBadRequest)
		return
	}
	log.Println("Summary Response: ", summaryResponse)
	if err != nil && strings.Contains(err.Error(), "no sessions found for the given day") {
		HandleError(c, fmt.Errorf("no sessions found for user: %s on date: %s", userId, queryDate), http.StatusNoContent)
		return
	}

	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, summaryResponse)

}

func (h *Handler) GetWeeklySummary(c *gin.Context) {

	userId := strings.TrimSpace(c.GetHeader("X-USER-ID"))
	if userId == "" {
		HandleError(c, fmt.Errorf("user id is missing"), http.StatusUnauthorized)
		return
	}

	start := c.Query("start_date")

	summaryResponse, err := h.Svc.GetWeeklySessionSummary(userId, start)

	if err != nil && strings.Contains(err.Error(), "invalid date format") {
		HandleError(c, err, http.StatusBadRequest)
		return
	}

	if err != nil && strings.Contains(err.Error(), "no sessions found") {
		HandleError(c, fmt.Errorf("no sessions found for user: %s", userId), http.StatusNoContent)
		return
	}

	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, summaryResponse)

}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "summary-service",
		"timestamp": time.Now(),
		"profile":   "production",
	})
}

func (h *Handler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
