package ratelimiter

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestMiddleware_Disabled(t *testing.T) {
	config := &Config{
		Enabled:         false,
		Capacity:        5.0,
		RefillRate:      1.0,
		CleanupInterval: 1 * time.Minute,
		TTL:             30 * time.Minute,
	}

	manager := NewRateLimiterManager(config.Capacity, config.RefillRate, config.CleanupInterval, config.TTL)
	defer manager.Shutdown()

	middleware := Middleware(manager, config)

	// Create a simple handler that returns 200 OK
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Wrap with middleware
	wrappedHandler := middleware(handler)

	// Make multiple requests - all should pass since rate limiting is disabled
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"

		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d (rate limiting should be disabled)", i+1, w.Code)
		}
	}
}

func TestMiddleware_AllowedRequest(t *testing.T) {
	config := &Config{
		Enabled:         true,
		Capacity:        5.0,
		RefillRate:      1.0,
		CleanupInterval: 1 * time.Minute,
		TTL:             30 * time.Minute,
	}

	manager := NewRateLimiterManager(config.Capacity, config.RefillRate, config.CleanupInterval, config.TTL)
	defer manager.Shutdown()

	middleware := Middleware(manager, config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	// Should be allowed
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check rate limit headers
	limitHeader := w.Header().Get("X-RateLimit-Limit")
	if limitHeader != "5" {
		t.Errorf("Expected X-RateLimit-Limit: 5, got %s", limitHeader)
	}

	remainingHeader := w.Header().Get("X-RateLimit-Remaining")
	if remainingHeader != "4" {
		t.Errorf("Expected X-RateLimit-Remaining: 4, got %s", remainingHeader)
	}
}

func TestMiddleware_RateLimitExceeded(t *testing.T) {
	config := &Config{
		Enabled:         true,
		Capacity:        3.0,
		RefillRate:      1.0,
		CleanupInterval: 1 * time.Minute,
		TTL:             30 * time.Minute,
	}

	manager := NewRateLimiterManager(config.Capacity, config.RefillRate, config.CleanupInterval, config.TTL)
	defer manager.Shutdown()

	middleware := Middleware(manager, config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	wrappedHandler := middleware(handler)

	clientIP := "192.168.1.1:12345"

	// Make requests up to capacity
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = clientIP

		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d", i+1, w.Code)
		}
	}

	// Next request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = clientIP

	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w.Code)
	}

	// Check response body
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["error"] != "rate limit exceeded" {
		t.Errorf("Expected error message, got %v", response)
	}

	// Check headers
	remainingHeader := w.Header().Get("X-RateLimit-Remaining")
	if remainingHeader != "0" {
		t.Errorf("Expected X-RateLimit-Remaining: 0, got %s", remainingHeader)
	}
}

func TestMiddleware_DifferentIPsIndependent(t *testing.T) {
	config := &Config{
		Enabled:         true,
		Capacity:        2.0,
		RefillRate:      1.0,
		CleanupInterval: 1 * time.Minute,
		TTL:             30 * time.Minute,
	}

	manager := NewRateLimiterManager(config.Capacity, config.RefillRate, config.CleanupInterval, config.TTL)
	defer manager.Shutdown()

	middleware := Middleware(manager, config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middleware(handler)

	// Exhaust IP1's rate limit
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"

		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)
	}

	// IP1 should be rate limited
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	w1 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusTooManyRequests {
		t.Errorf("IP1 should be rate limited, got status %d", w1.Code)
	}

	// IP2 should still be allowed
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.2:12345"
	w2 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("IP2 should be allowed, got status %d", w2.Code)
	}
}

func TestMiddleware_RefillAllowsNewRequests(t *testing.T) {
	config := &Config{
		Enabled:         true,
		Capacity:        2.0,
		RefillRate:      10.0, // 10 tokens per second for faster test
		CleanupInterval: 1 * time.Minute,
		TTL:             30 * time.Minute,
	}

	manager := NewRateLimiterManager(config.Capacity, config.RefillRate, config.CleanupInterval, config.TTL)
	defer manager.Shutdown()

	middleware := Middleware(manager, config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middleware(handler)

	clientIP := "192.168.1.1:12345"

	// Exhaust rate limit
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = clientIP
		w := httptest.NewRecorder()
		wrappedHandler.ServeHTTP(w, req)
	}

	// Should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = clientIP
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected rate limit, got status %d", w.Code)
	}

	// Wait for refill (100ms should add 1 token at 10 tokens/sec)
	time.Sleep(100 * time.Millisecond)

	// Should now be allowed
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = clientIP
	w = httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected request to be allowed after refill, got status %d", w.Code)
	}
}

func TestExtractClientIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.100:54321"

	ip := extractClientIP(req)

	if ip != "192.168.1.100" {
		t.Errorf("Expected IP 192.168.1.100, got %s", ip)
	}
}

func TestExtractClientIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 198.51.100.1")

	ip := extractClientIP(req)

	// Should use first IP from X-Forwarded-For
	if ip != "203.0.113.1" {
		t.Errorf("Expected IP 203.0.113.1, got %s", ip)
	}
}

func TestExtractClientIP_XRealIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	req.Header.Set("X-Real-IP", "203.0.113.5")

	ip := extractClientIP(req)

	// Should use X-Real-IP
	if ip != "203.0.113.5" {
		t.Errorf("Expected IP 203.0.113.5, got %s", ip)
	}
}

func TestExtractClientIP_XForwardedForPriority(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.1")
	req.Header.Set("X-Real-IP", "203.0.113.5")

	ip := extractClientIP(req)

	// X-Forwarded-For should take priority
	if ip != "203.0.113.1" {
		t.Errorf("Expected X-Forwarded-For to take priority, got %s", ip)
	}
}

func TestMiddleware_HandlerNotCalledWhenRateLimited(t *testing.T) {
	config := &Config{
		Enabled:         true,
		Capacity:        1.0,
		RefillRate:      1.0,
		CleanupInterval: 1 * time.Minute,
		TTL:             30 * time.Minute,
	}

	manager := NewRateLimiterManager(config.Capacity, config.RefillRate, config.CleanupInterval, config.TTL)
	defer manager.Shutdown()

	middleware := Middleware(manager, config)

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middleware(handler)

	clientIP := "192.168.1.1:12345"

	// First request - should call handler
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = clientIP
	w1 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w1, req1)

	if !handlerCalled {
		t.Error("Handler should have been called for first request")
	}

	// Reset flag
	handlerCalled = false

	// Second request - should NOT call handler (rate limited)
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = clientIP
	w2 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w2, req2)

	if handlerCalled {
		t.Error("Handler should NOT have been called when rate limited")
	}

	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", w2.Code)
	}
}
