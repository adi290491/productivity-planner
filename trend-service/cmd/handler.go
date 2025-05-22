package main

import (
	"fmt"
	"log"
	"net/http"
	"productivity-planner/trend-service/trend"
	"productivity-planner/trend-service/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *trend.TrendService
}

func NewHandler(svc *trend.TrendService) *Handler {
	return &Handler{
		svc: svc,
	}
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

	weelyTrendResponse, err := h.svc.FetchWeeklyTrend(userId, weeks)

	if err != nil {
		HandleError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, weelyTrendResponse)

}
