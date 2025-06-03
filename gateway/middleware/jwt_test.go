package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestJWTMiddleware(t *testing.T) {
	secret := "mysecret"
	validToken := generateValidJWT(secret, "1234")

	tests := []struct {
		name               string
		token              string
		expectedStatusCode int
	}{
		{
			name:               "Valid token",
			token:              "Bearer " + validToken,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Missing token",
			token:              "",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Invalid token format",
			token:              "Token abcdef",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Invalid token content",
			token:              "Bearer invalid.token.here",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("JWT_SECRET", secret)

			r := gin.New()
			r.Use(JWTMiddleware())
			r.GET("/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "authorized"})
			})

			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, w.Code)
			}
		})
	}
}

func generateValidJWT(secret, userID string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}
