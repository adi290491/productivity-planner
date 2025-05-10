package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing auth token"})
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		jwtUtil := JWTUtil{
			Secret:      os.Getenv("JWT_SECRET"),
			tokenString: tokenStr,
		}

		userId, err := jwtUtil.ValidateToken()

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
