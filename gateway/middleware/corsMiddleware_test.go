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

	// Test with a valid origin
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://systemic-productivity-planner.web.app")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	headers := w.Header()
	tests := map[string]string{
		"Access-Control-Allow-Origin":      "https://systemic-productivity-planner.web.app",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization, X-USER-ID",
		"Access-Control-Allow-Credentials": "true",
	}

	for key, expected := range tests {
		if val := headers.Get(key); val != expected {
			t.Errorf("Expected header %s to be %s, got %s", key, expected, val)
		}
	}

	// Test with an invalid origin (should not get CORS headers)
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.Header.Set("Origin", "https://malicious-site.com")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	// Should not have Access-Control-Allow-Origin for invalid origin
	if origin := w2.Header().Get("Access-Control-Allow-Origin"); origin != "" {
		t.Errorf("Expected no Access-Control-Allow-Origin for invalid origin, got %s", origin)
	}
}
