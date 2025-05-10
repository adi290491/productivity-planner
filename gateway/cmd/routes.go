package main

import (
	"productivity-planner/gateway/middleware"
	"productivity-planner/gateway/proxy"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	{
		usersRouter := r.Group("/")
		usersRouter.POST("/users/signup", proxy.ProxyToUserService)
		usersRouter.POST("/users/login", proxy.ProxyToUserService)
	}

	{
		sessionsRouter := r.Group("/sessions")
		sessionsRouter.Use(middleware.JWTMiddleware())
		sessionsRouter.POST("/v1/start-session", proxy.ProxyToSessionService)
		sessionsRouter.PATCH("/v1/stop-session", proxy.ProxyToSessionService)
	}

	{
		summaryRouter := r.Group("/summary")
		summaryRouter.Use(middleware.JWTMiddleware())
		summaryRouter.GET("/daily", proxy.ProxyToSummaryService)
		summaryRouter.GET("/weekly", proxy.ProxyToSummaryService)
	}
}
