package main

import (
	"github.com/adi290491/productivity-planner/gateway/config"
	"github.com/adi290491/productivity-planner/gateway/middleware"
	"github.com/adi290491/productivity-planner/gateway/proxy"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, cfg *config.AppConfig) {

	// Use only one CORS middleware to avoid conflicts
	r.Use(middleware.CorsMiddleware())
	{
		usersRouter := r.Group("/")
		usersRouter.POST("/users/signup", proxy.ProxyToUserService)
		usersRouter.POST("/users/login", proxy.ProxyToUserService)
	}

	{
		sessionsRouter := r.Group("/sessions")
		sessionsRouter.Use(middleware.JWTMiddleware(cfg))
		sessionsRouter.POST("/v1/start-session", proxy.ProxyToSessionService)
		sessionsRouter.PATCH("/v1/stop-session", proxy.ProxyToSessionService)
	}

	{
		summaryRouter := r.Group("/summary")
		summaryRouter.Use(middleware.JWTMiddleware(cfg))
		summaryRouter.GET("/daily", proxy.ProxyToSummaryService)
		summaryRouter.GET("/weekly", proxy.ProxyToSummaryService)
	}

	{
		trendRouter := r.Group("/trend")
		trendRouter.Use(middleware.JWTMiddleware(cfg))
		trendRouter.GET("/daily", proxy.ProxyToTrendService)
		trendRouter.GET("/weekly", proxy.ProxyToTrendService)
		trendRouter.GET("/unviewed", proxy.ProxyToTrendService)
		trendRouter.POST("/mark-viewed", proxy.ProxyToTrendService)
	}

}
