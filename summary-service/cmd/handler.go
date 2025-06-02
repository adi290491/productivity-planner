package main

import (
	"fmt"
	"log"
	"net/http"
	"productivity-planner/summary-service/summary"
	"strings"

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
		HandleError(c, fmt.Errorf("no sessions found for user: %s on date: %s", userId, queryDate), http.StatusNotFound)
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
		HandleError(c, fmt.Errorf("no sessions found for user: %s", userId), http.StatusNotFound)
		return
	}

	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, summaryResponse)

}
