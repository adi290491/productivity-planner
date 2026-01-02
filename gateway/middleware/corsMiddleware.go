package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Define allowed origins (fallback for local development)
		allowedOrigins := []string{
			"https://systemic-productivity-planner.web.app",
			"https://systemic-productivity-planner.firebaseapp.com",
			"http://localhost:4200",
			"http://localhost:5173",
			"http://localhost:3000",
		}

		// Add environment-based origins (production)
		if frontendOrigin := os.Getenv("FRONTEND_ORIGIN"); frontendOrigin != "" {
			allowedOrigins = append([]string{frontendOrigin}, allowedOrigins...)
		}

		if frontendOrigin2 := os.Getenv("FRONTEND_ORIGIN_2"); frontendOrigin2 != "" {
			allowedOrigins = append([]string{frontendOrigin2}, allowedOrigins...)
		}

		// Check if origin is in allowed list
		originAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				// c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				originAllowed = true
				break
			}
		}

		// Always set other CORS headers for valid origins
		if originAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-USER-ID")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			if originAllowed {
				c.AbortWithStatus(http.StatusNoContent)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}
		c.Next()
	}
}
