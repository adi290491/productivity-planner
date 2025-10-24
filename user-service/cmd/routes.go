package main

import (
	"github.com/gin-gonic/gin"
)

func RegisterEndpoints(r *gin.Engine, h *Handler) {

	r.GET("/health", h.HealthCheck)
	r.GET("/ready", h.Ready)

	r.POST("/users/signup", h.Signup)
	r.POST("/users/login", h.Login)
	r.POST("/users/batch", h.GetUsersBatch)
}
