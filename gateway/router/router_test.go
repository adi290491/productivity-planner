package router

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/adi290491/productivity-planner/gateway/config"
	"github.com/adi290491/productivity-planner/gateway/ratelimiter"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	// Suppress logs during tests
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func setupTestRouter() http.Handler {
	// Set up test config
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("USER_SERVICE_URL", "http://user-service:8081")
	os.Setenv("SESSION_SERVICE_URL", "http://session-service:8085")
	os.Setenv("SUMMARY_SERVICE_URL", "http://summary-service:8082")
	os.Setenv("TREND_SERVICE_URL", "http://trend-service:8083")

	cfg := &config.AppConfig{
		JWT_SECRET:   "test-secret",
		Port:         "8000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	rateLimitCfg := &ratelimiter.Config{
		Enabled:         true,
		Capacity:        10.0,
		RefillRate:      5.0,
		CleanupInterval: 1 * time.Minute,
		TTL:             30 * time.Minute,
	}

	rateLimiterMgr := ratelimiter.NewRateLimiterManager(
		rateLimitCfg.Capacity,
		rateLimitCfg.RefillRate,
		rateLimitCfg.CleanupInterval,
		rateLimitCfg.TTL,
	)

	return NewRouter(cfg, rateLimiterMgr, rateLimitCfg)
}

func generateTestToken(secret, userId string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func TestRouter_PublicRoute_NoAuth(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest("POST", "/users/signup", nil)
	req.Header.Set("Origin", "http://localhost:5173")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should succeed without auth token
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check it's the user service placeholder
	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["message"] != "user-service - placeholder" {
		t.Errorf("Expected user-service response, got %v", response)
	}
}

func TestRouter_ProtectedRoute_WithAuth(t *testing.T) {
	router := setupTestRouter()

	token := generateTestToken("test-secret", "user123")

	req := httptest.NewRequest("POST", "/sessions/v1/start-session", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should succeed with valid token
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRouter_ProtectedRoute_WithoutAuth(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest("POST", "/sessions/v1/start-session", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	// No Authorization header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should fail without token
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestRouter_CORS_Preflight(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest("OPTIONS", "/users/signup", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 204 No Content for preflight
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Should have CORS headers
	if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:5173" {
		t.Errorf("Expected CORS header, got %s", origin)
	}
}

func TestRouter_RateLimiting(t *testing.T) {
	router := setupTestRouter()

	// Make requests up to capacity (10)
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("POST", "/users/signup", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		req.RemoteAddr = "192.168.1.100:12345"

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// 11th request should be rate limited
	req := httptest.NewRequest("POST", "/users/signup", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.RemoteAddr = "192.168.1.100:12345"

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w.Code)
	}
}

func TestRouter_AllRoutes(t *testing.T) {
	router := setupTestRouter()
	token := generateTestToken("test-secret", "user123")

	tests := []struct {
		name           string
		path           string
		needsAuth      bool
		expectedStatus int
	}{
		// Public routes
		{"User signup", "/users/signup", false, http.StatusOK},
		{"User login", "/users/login", false, http.StatusOK},

		// Protected routes
		{"Start session", "/sessions/v1/start-session", true, http.StatusOK},
		{"Stop session", "/sessions/v1/stop-session", true, http.StatusOK},
		{"Daily summary", "/summary/daily", true, http.StatusOK},
		{"Weekly summary", "/summary/weekly", true, http.StatusOK},
		{"Daily trend", "/trend/daily", true, http.StatusOK},
		{"Weekly trend", "/trend/weekly", true, http.StatusOK},
		{"Unviewed trends", "/trend/unviewed", true, http.StatusOK},
		{"Mark viewed", "/trend/mark-viewed", true, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", tt.path, nil)
			req.Header.Set("Origin", "http://localhost:5173")
			req.RemoteAddr = "192.168.1.50:12345" // Different IP for each test

			if tt.needsAuth {
				req.Header.Set("Authorization", "Bearer "+token)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRouter_ProtectedRoute_InvalidToken(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest("POST", "/sessions/v1/start-session", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestRouter_ProtectedRoute_ExpiredToken(t *testing.T) {
	router := setupTestRouter()

	// Create expired token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": "user123",
		"exp":    time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	})
	tokenStr, _ := token.SignedString([]byte("test-secret"))

	req := httptest.NewRequest("POST", "/sessions/v1/start-session", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestRouter_DifferentIPsIndependent(t *testing.T) {
	router := setupTestRouter()

	// Exhaust rate limit for IP1
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("POST", "/users/signup", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		req.RemoteAddr = "192.168.1.1:12345"

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// IP1 should be rate limited
	req1 := httptest.NewRequest("POST", "/users/signup", nil)
	req1.Header.Set("Origin", "http://localhost:5173")
	req1.RemoteAddr = "192.168.1.1:12345"

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusTooManyRequests {
		t.Errorf("IP1 should be rate limited, got %d", w1.Code)
	}

	// IP2 should still work
	req2 := httptest.NewRequest("POST", "/users/signup", nil)
	req2.Header.Set("Origin", "http://localhost:5173")
	req2.RemoteAddr = "192.168.1.2:12345"

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("IP2 should work, got %d", w2.Code)
	}
}
