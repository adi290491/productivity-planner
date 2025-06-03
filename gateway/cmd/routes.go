package main

import (
	"os"
	"productivity-planner/gateway/middleware"
	"productivity-planner/gateway/proxy"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")

	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:5173" // fallback
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-USER-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.CorsMiddleware())
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

	{
		trendRouter := r.Group("/trend")
		trendRouter.Use(middleware.JWTMiddleware())
		trendRouter.GET("/daily", proxy.ProxyToTrendService)
		trendRouter.GET("/weekly", proxy.ProxyToTrendService)
	}
}
