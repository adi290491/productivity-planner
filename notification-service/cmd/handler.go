package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/adi290491/productivity-planner/notification-service/config"
	"github.com/adi290491/productivity-planner/notification-service/models"
	"github.com/adi290491/productivity-planner/notification-service/notification"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NotificationServiceInterface defines the interface for notification service
type NotificationServiceInterface interface {
	ProcessDailyTrendNotifications(ctx context.Context) error
	ProcessWeeklyTrendNotifications(ctx context.Context) error
	GetUserNotification(userID uuid.UUID) (*models.UserNotificationResponse, error)
	MarkNotificationAsRead(userID uuid.UUID, notificationType string) error
	GetStats() notification.ProcessingStats
}

type Handler struct {
	config              *config.AppConfig
	notificationService NotificationServiceInterface
}

func NewHandler(config *config.AppConfig, notificationService NotificationServiceInterface) *Handler {
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
	// Extract authenticated user ID from request context (set by auth middleware)
	authenticatedUserIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	authenticatedUserID, ok := authenticatedUserIDValue.(uuid.UUID)
	if !ok {
		// Try to parse if it's a string
		if userIDStr, isString := authenticatedUserIDValue.(string); isString {
			var err error
			authenticatedUserID, err = uuid.Parse(userIDStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authenticated user ID"})
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authenticated user ID"})
			return
		}
	}

	// Get user_id from query parameter
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	requestedUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	// Verify that authenticated user can only access their own notifications
	if authenticatedUserID != requestedUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: can only access your own notifications"})
		return
	}

	response, err := h.notificationService.GetUserNotification(authenticatedUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// markNotificationAsRead marks specific notification types as read for a user
func (h *Handler) markNotificationAsRead(c *gin.Context) {
	// Extract authenticated user ID from request context (set by auth middleware)
	authenticatedUserIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	authenticatedUserID, ok := authenticatedUserIDValue.(uuid.UUID)
	if !ok {
		// Try to parse if it's a string
		if userIDStr, isString := authenticatedUserIDValue.(string); isString {
			var err error
			authenticatedUserID, err = uuid.Parse(userIDStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authenticated user ID"})
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authenticated user ID"})
			return
		}
	}

	// Get user_id from query parameter
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	requestedUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	// Verify that authenticated user can only access their own notifications
	if authenticatedUserID != requestedUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: can only mark your own notifications as read"})
		return
	}

	// Validate and normalize notification type
	notificationType := strings.ToLower(strings.TrimSpace(c.Query("type")))
	if notificationType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "notification type is required"})
		return
	}

	// Validate notification type accepts only "daily" or "weekly"
	if notificationType != "daily" && notificationType != "weekly" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid notification type: must be 'daily' or 'weekly'",
		})
		return
	}

	err = h.notificationService.MarkNotificationAsRead(authenticatedUserID, notificationType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("%s notification marked as read", notificationType),
	})
}
