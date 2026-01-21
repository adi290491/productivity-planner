package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/adi290491/productivity-planner/user-service/user"
	"github.com/adi290491/productivity-planner/user-service/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Svc     user.UserServiceInterface
	JwtUtil utils.JWTUtil
}

var (
	dbInitialized atomic.Bool
	dbInitMutex   sync.Mutex
)

func (h *Handler) Signup(c *gin.Context) {
	slog.Info("Signup request received")

	var req user.SignupDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.Svc.Signup(req)

	if err != nil {
		slog.Error("Signup failed",
			"method", c.Request.Method,
			"path", c.FullPath(),
			"ip", c.ClientIP(),
			"error", err,
		)
		HandleError(c, err, 500)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func (h *Handler) Login(c *gin.Context) {
	slog.Info("Login request received")

	var req user.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to parse login request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	slog.Debug("Attempting login", "email", req.Email)

	user, err := h.Svc.Login(req)

	if err != nil {
		slog.Error("Login service failed",
			"email", req.Email,
			"error", err,
		)
		HandleError(c, err, 400)
		return
	}

	slog.Debug("Login successful, generating token", "userId", user.ID)

	token, err := h.JwtUtil.GenerateToken(user)

	if err != nil {
		slog.Error("Token generation failed",
			"userId", user.ID,
			"error", err,
		)
		HandleError(c, err, 500)
		return
	}

	slog.Info("Login completed successfully", "userId", user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) GetUsersBatch(c *gin.Context) {
	var req user.GetUsersBatchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, fmt.Errorf("invalid input"), http.StatusBadRequest)
		return
	}

	if len(req.UserIDs) == 0 {
		c.JSON(http.StatusBadRequest, []user.UserInfoResponse{})
		return
	}

	response, err := h.Svc.GetUsersBatch(req)

	if err != nil {
		HandleError(c, fmt.Errorf("failed to get users batch: %w", err), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	if dbInitialized.Load() {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "user-service",
			"timestamp": time.Now(),
			"profile":   "production",
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "db_not_ready",
			"service":   "user-service",
			"timestamp": time.Now(),
			"profile":   "production",
		})
	}
}

func (h *Handler) Ready(c *gin.Context) {
	if dbInitialized.Load() {
		c.JSON(http.StatusOK, gin.H{"status": "ready", "db": "connected"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "starting", "db": "initializing"})
	}
}
