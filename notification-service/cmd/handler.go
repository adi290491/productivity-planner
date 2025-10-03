package main

import (
	"fmt"
	"net/http"

	"github.com/adi290491/productivity-planner/notification-service/config"
	"github.com/adi290491/productivity-planner/notification-service/notification"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	config              *config.AppConfig
	notificationService *notification.NotificationService
}

func NewHandler(config *config.AppConfig, notificationService *notification.NotificationService) *Handler {
	return &Handler{
		config:              config,
		notificationService: notificationService,
	}
}

func (h *Handler) processDailyTrends(c *gin.Context) {
	ctx := c.Request.Context()

	err := h.notificationService.ProcessDailyTrendNotifications(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	stats := h.notificationService.GetStats()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daily trend notifications processed successfully",
		"stats":   stats,
	})
}

func (h *Handler) processWeeklyTrends(c *gin.Context) {
	ctx := c.Request.Context()

	err := h.notificationService.ProcessWeeklyTrendNotifications(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	stats := h.notificationService.GetStats()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Weekly trend notifications processed successfully",
		"stats":   stats,
	})
}

func (h *Handler) getUserNotifications(c *gin.Context) {
	// TODO: Get user ID from JWT token in production
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	response, err := h.notificationService.GetUserNotification(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// markNotificationAsRead marks specific notification types as read for a user
func (h *Handler) markNotificationAsRead(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	notificationType := c.Query("type") // "daily" or "weekly"
	if notificationType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "notification type is required (daily or weekly)"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	err = h.notificationService.MarkNotificationAsRead(userID, notificationType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("%s notification marked as read", notificationType),
	})
}
