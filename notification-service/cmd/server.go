package main

import (
	"log"
	"net/http"
	"time"

	"github.com/adi290491/productivity-planner/notification-service/config"
	"github.com/adi290491/productivity-planner/notification-service/notification"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config              *config.AppConfig
	notificationService *notification.NotificationService
	handler             *Handler
	router              *gin.Engine
}

func NewServer(config *config.AppConfig, notificationService *notification.NotificationService) *Server {
	if config.Profile == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	handler := NewHandler(config, notificationService)

	server := &Server{
		config:              config,
		notificationService: notificationService,
		handler:             handler,
		router:              gin.Default(),
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthCheck)

	// Statistics
	s.router.GET("/stat", s.getStats)

	// Configuration endpoint
	s.router.GET("/config", s.getConfig)

	// Handler routes
	s.router.POST("/process/daily", s.handler.processDailyTrends)
	s.router.POST("/process/weekly", s.handler.processWeeklyTrends)
	s.router.GET("/notifications", s.handler.getUserNotifications)
	s.router.POST("/notifications/:userID/read", s.handler.markNotificationAsRead)
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "notification-service",
		"timestamp": time.Now(),
		"profile":   s.config.Profile,
	})

}

func (s *Server) getStats(c *gin.Context) {
	stats := s.notificationService.GetStats()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Daily trend notifications processed successfully",
		"stats":   stats,
	})
}

func (s *Server) getConfig(c *gin.Context) {
	// Return safe config (without sensitive information)
	safeConfig := map[string]interface{}{
		"profile":            s.config.Profile,
		"port":               s.config.Port,
		"project_id":         s.config.ProjectID,
		"daily_topic":        s.config.DailyTopic,
		"daily_subscription": s.config.DailySubscription,
		// "admin_email":       s.config.AdminEmail,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"config":  safeConfig,
	})
}

func (s *Server) Run() error {
	log.Printf("Starting notification service on port %s", s.config.Port)
	return s.router.Run(":" + s.config.Port)
}
