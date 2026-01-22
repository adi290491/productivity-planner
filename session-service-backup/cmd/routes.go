package main

import "github.com/gin-gonic/gin"

func RegisterEndpoints(r *gin.Engine, h *Handler) {

	r.GET("/health", h.HealthCheck)
	r.GET("/ready", h.Ready)

	r.POST("/sessions/v1/start-session", h.StartSession)
	r.PATCH("/sessions/v1/stop-session", h.StopSession)
}
