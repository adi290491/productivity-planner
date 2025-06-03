package proxy

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestProxyToUserService(t *testing.T) {
	mockServer := startMockBackend(t)
	defer mockServer.Close()

	os.Setenv("USER_SERVICE_URL", mockServer.URL)

	router := gin.New()
	router.POST("/users/signup", ProxyToUserService)

	body := []byte(`{"email":"test@example.com","password":"pass","name":"Test User"}`)
	req, _ := http.NewRequest("POST", "/users/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestProxyToSessionService(t *testing.T) {
	mockServer := startMockBackend(t)
	defer mockServer.Close()

	os.Setenv("SESSION_SERVICE_URL", mockServer.URL)

	router := gin.New()
	router.POST("/sessions/v1/start-session", func(c *gin.Context) {
		// Add fake userId to context
		c.Set("userId", "mock-user-id")
		ProxyToSessionService(c)
	})

	body := []byte(`{"session_type":"focus"}`)
	req, _ := http.NewRequest("POST", "/sessions/v1/start-session", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Header().Get("X-USER-ID") == "" {
		t.Errorf("Expected X-USER-ID header to be set")
	}
}

func TestProxyToSummaryService(t *testing.T) {
	mockServer := startMockBackend(t)
	defer mockServer.Close()

	os.Setenv("SUMMARY_SERVICE_URL", mockServer.URL)

	router := gin.New()
	router.GET("/summary/daily", func(c *gin.Context) {
		c.Set("userId", "mock-user-id")
		ProxyToSummaryService(c)
	})

	req, _ := http.NewRequest("GET", "/summary/daily?date=2025-05-01", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestProxyToTrendService(t *testing.T) {
	mockServer := startMockBackend(t)
	defer mockServer.Close()

	os.Setenv("TREND_SERVICE_URL", mockServer.URL)

	router := gin.New()
	router.GET("/trend/weekly", func(c *gin.Context) {
		c.Set("userId", "mock-user-id")
		ProxyToTrendService(c)
	})

	req, _ := http.NewRequest("GET", "/trend/weekly?weeks=2", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func startMockBackend(t *testing.T) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-USER-ID")
		w.Header().Set("X-USER-ID", userID) // Echo it back for testing

		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body) // Just echoing the request body back
	})
	return httptest.NewServer(handler)
}
