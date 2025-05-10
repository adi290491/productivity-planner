package main

import (
	"github.com/gin-gonic/gin"
)

func RegisterEndpoints(r *gin.Engine, h *Handler) {

	r.POST("/users/signup", h.Signup)
	r.POST("/users/login", h.Login)
}
