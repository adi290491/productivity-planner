package middleware

import (
	"net/http"
	"strings"

	"github.com/adi290491/productivity-planner/gateway/config"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(cfg *config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid auth token"})
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Empty auth token"})
			return
		}

		jwtUtil := JWTUtil{
			Secret:      cfg.JWT_SECRET,
			tokenString: tokenStr,
		}

		userId, err := jwtUtil.ValidateToken()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if userId == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
