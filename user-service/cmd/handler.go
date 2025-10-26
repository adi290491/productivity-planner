package main

import (
	"fmt"
	"log"
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
	log.Println("--------Called Signup function---------")
	var req user.SignupDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.Svc.Signup(req)

	if err != nil {
		log.Println("--------Error on Signup---------")
		log.Printf("Signup failed method=%s path=%s ip=%s err=%v", c.Request.Method, c.FullPath(), c.ClientIP(), err)
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
	var req user.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.Svc.Login(req)

	if err != nil {
		HandleError(c, err, 400)
		return
	}

	token, err := h.JwtUtil.GenerateToken(user)

	if err != nil {
		HandleError(c, err, 500)
		return
	}

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
