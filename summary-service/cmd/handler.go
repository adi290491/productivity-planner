package main

import (
	"fmt"
	"net/http"
	"productivity-planner/summary-service/summary"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *summary.SummaryService
}

func NewHandler(svc *summary.SummaryService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) GetDailySummary(c *gin.Context) {

	userId := strings.TrimSpace(c.GetHeader("X-USER-ID"))
	if userId == "" {
		HandleError(c, fmt.Errorf("user id is missing"), http.StatusUnauthorized)
		return
	}

	queryDate := c.Query("date")
	var day time.Time
	var err error

	if queryDate == "" {
		day = time.Now().UTC()
	} else {
		day, err = time.Parse("2006-01-02", queryDate)
		if err != nil {
			HandleError(c, fmt.Errorf("invalid date format"), http.StatusBadRequest)
			return
		}
	}

	summaryResponse, err := h.svc.GetDailySessionSummary(userId, day)

	if err != nil && strings.Contains(err.Error(), "no sessions found for the given day") {
		HandleError(c, fmt.Errorf("no sessions found for user: %s on date: %s", userId, day.Format("2006-01-02")), http.StatusNotFound)
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

	summaryResponse, err := h.svc.GetWeeklySessionSummary(userId, start)

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
