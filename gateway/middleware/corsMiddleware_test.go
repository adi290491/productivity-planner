package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCorsMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(CorsMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	headers := w.Header()
	tests := map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization",
		"Access-Control-Allow-Credentials": "true",
	}

	for key, expected := range tests {
		if val := headers.Get(key); val != expected {
			t.Errorf("Expected header %s to be %s, got %s", key, expected, val)
		}
	}
}
