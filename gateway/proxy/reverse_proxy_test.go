package proxy

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
	"time"
)

func init() {
	// Suppress logs during tests
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

// createTestProxy creates a proxy with mocked authentication for testing
func createTestProxy(targetURL, service string) (http.Handler, error) {
	proxy, err := NewReverseProxy(ProxyConfig{
		TargetURL: targetURL,
		Service:   service,
	})
	if err != nil {
		return nil, err
	}

	// Mock the token fetcher to avoid needing real Google Cloud credentials
	// This is safe because httptest servers don't check authentication
	if reverseProxy, ok := proxy.(*httputil.ReverseProxy); ok {
		if transport, ok := reverseProxy.Transport.(*AuthenticatedTransport); ok {
			transport.tokenFetcher = func(ctx context.Context, audience string) (string, time.Time, error) {
				return "test-token", time.Now().Add(1 * time.Hour), nil
			}
		}
	}

	return proxy, nil
}

func TestNewReverseProxy_ValidURL(t *testing.T) {
	// Create a mock backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer backend.Close()

	proxy, err := createTestProxy(backend.URL, "test-service")

	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	if proxy == nil {
		t.Fatal("Expected non-nil proxy")
	}
}

func TestNewReverseProxy_InvalidURL(t *testing.T) {
	_, err := NewReverseProxy(ProxyConfig{
		TargetURL: "://invalid-url",
		Service:   "test-service",
	})

	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestNewReverseProxy_ForwardsRequest(t *testing.T) {
	// Track what the backend receives
	var receivedMethod string
	var receivedPath string
	var receivedHeaders http.Header

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedPath = r.URL.Path
		receivedHeaders = r.Header.Clone()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"backend response"}`))
	}))
	defer backend.Close()

	proxy, err := createTestProxy(backend.URL, "test-service")
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	// Make request through proxy
	req := httptest.NewRequest("POST", "/api/test", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	proxy.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check backend received correct method and path
	if receivedMethod != "POST" {
		t.Errorf("Expected method POST, backend received %s", receivedMethod)
	}

	if receivedPath != "/api/test" {
		t.Errorf("Expected path /api/test, backend received %s", receivedPath)
	}

	// Check Content-Type was forwarded
	if ct := receivedHeaders.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, backend received %s", ct)
	}
}

func TestNewReverseProxy_AddsAuthorizationHeader(t *testing.T) {
	var receivedAuth string

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy, err := createTestProxy(backend.URL, "test-service")
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	proxy.ServeHTTP(w, req)

	// Check Authorization header was added
	if receivedAuth != "Bearer test-token" {
		t.Errorf("Expected Authorization: Bearer test-token, got %s", receivedAuth)
	}
}

func TestNewReverseProxy_AddsUserIDHeader(t *testing.T) {
	var receivedUserID string

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUserID = r.Header.Get("X-User-ID")
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy, err := createTestProxy(backend.URL, "test-service")
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	// Create request with userId in context
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), "userId", "user123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	proxy.ServeHTTP(w, req)

	// Check X-User-ID was added
	if receivedUserID != "user123" {
		t.Errorf("Expected X-User-ID: user123, backend received %s", receivedUserID)
	}
}

func TestNewReverseProxy_NoUserIDInContext(t *testing.T) {
	var receivedUserID string

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUserID = r.Header.Get("X-User-ID")
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy, err := createTestProxy(backend.URL, "test-service")
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	// Create request WITHOUT userId in context
	req := httptest.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	proxy.ServeHTTP(w, req)

	// X-User-ID should not be set
	if receivedUserID != "" {
		t.Errorf("Expected no X-User-ID header, but got %s", receivedUserID)
	}
}

func TestNewReverseProxy_StripsCORSHeaders(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Backend sets CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy, err := createTestProxy(backend.URL, "test-service")
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	proxy.ServeHTTP(w, req)

	// CORS headers should be stripped
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("Expected CORS headers to be stripped")
	}
}

func TestNewReverseProxy_ErrorHandler(t *testing.T) {
	// Create a proxy pointing to a non-existent backend
	proxy, err := createTestProxy("http://localhost:99999", "test-service")
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	proxy.ServeHTTP(w, req)

	// Should return 502 Bad Gateway
	if w.Code != http.StatusBadGateway {
		t.Errorf("Expected status 502, got %d", w.Code)
	}

	// Check error response
	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["error"] != "service unavailable" {
		t.Errorf("Expected error message, got %v", response)
	}
}

func TestNewUserServiceProxy(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy, err := NewUserServiceProxy(backend.URL)
	if err != nil {
		t.Fatalf("Failed to create user service proxy: %v", err)
	}

	if proxy == nil {
		t.Fatal("Expected non-nil proxy")
	}

	// Mock the token fetcher
	if reverseProxy, ok := proxy.(*httputil.ReverseProxy); ok {
		if transport, ok := reverseProxy.Transport.(*AuthenticatedTransport); ok {
			transport.tokenFetcher = func(ctx context.Context, audience string) (string, time.Time, error) {
				return "test-token", time.Now().Add(1 * time.Hour), nil
			}
		}
	}

	// Test it works
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	proxy.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNewSessionServiceProxy(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy, err := NewSessionServiceProxy(backend.URL)
	if err != nil {
		t.Fatalf("Failed to create session service proxy: %v", err)
	}

	if proxy == nil {
		t.Fatal("Expected non-nil proxy")
	}
}

func TestNewSummaryServiceProxy(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy, err := NewSummaryServiceProxy(backend.URL)
	if err != nil {
		t.Fatalf("Failed to create summary service proxy: %v", err)
	}

	if proxy == nil {
		t.Fatal("Expected non-nil proxy")
	}
}

func TestNewTrendServiceProxy(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy, err := NewTrendServiceProxy(backend.URL)
	if err != nil {
		t.Fatalf("Failed to create trend service proxy: %v", err)
	}

	if proxy == nil {
		t.Fatal("Expected non-nil proxy")
	}
}

func TestGetUserIDFromContext_String(t *testing.T) {
	ctx := context.WithValue(context.Background(), "userId", "user123")

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		t.Error("Expected to find userID in context")
	}

	if userID != "user123" {
		t.Errorf("Expected userID 'user123', got %s", userID)
	}
}

func TestGetUserIDFromContext_NotPresent(t *testing.T) {
	ctx := context.Background()

	_, ok := GetUserIDFromContext(ctx)
	if ok {
		t.Error("Expected userID to not be in context")
	}
}

func TestGetUserIDFromContext_Interface(t *testing.T) {
	// Test with interface{} type (simulating what JWT middleware might set)
	var userID interface{} = "user456"
	ctx := context.WithValue(context.Background(), "userId", userID)

	result, ok := GetUserIDFromContext(ctx)
	if !ok {
		t.Error("Expected to find userID in context")
	}

	if result != "user456" {
		t.Errorf("Expected userID 'user456', got %s", result)
	}
}
