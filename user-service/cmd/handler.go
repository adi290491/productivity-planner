package main

import (
	"net/http"
	"productivity-planner/user-service/user"
	"productivity-planner/user-service/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Svc     *user.UserService
	JwtUtil utils.JWTUtil
}

func (h *Handler) Signup(c *gin.Context) {
	var req user.SignupDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.Svc.Signup(req)

	if err != nil {
		HandleError(c, err, 500)
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
	}

	token, err := h.JwtUtil.GenerateToken(user)

	if err != nil {
		HandleError(c, err, 500)
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
