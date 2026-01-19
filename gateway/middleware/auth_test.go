package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adi290491/productivity-planner/gateway/config"
	"github.com/golang-jwt/jwt/v5"
)

func generateTestToken(secret string, userId string, expiration time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    expiration.Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	secret := "test-secret"
	userId := "user123"

	cfg := &config.AppConfig{
		JWT_SECRET: secret,
	}

	middleware := JWTMiddleware(cfg)

	handlerCalled := false
	var contextUserId interface{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		contextUserId = r.Context().Value("userId")
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware(handler)

	token := generateTestToken(secret, userId, time.Now().Add(time.Hour))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if !handlerCalled {
		t.Error("Handler should have been called for valid token")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if contextUserId != userId {
		t.Errorf("Expected userId %s in context, got %v", userId, contextUserId)
	}
}

func TestJWTMiddleware_MissingToken(t *testing.T) {
	cfg := &config.AppConfig{
		JWT_SECRET: "test-secret",
	}

	middleware := JWTMiddleware(cfg)

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	// No Authorization header

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if handlerCalled {
		t.Error("Handler should not have been called without token")
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["error"] == "" {
		t.Error("Expected error message in response")
	}
}

func TestJWTMiddleware_InvalidTokenFormat(t *testing.T) {
	cfg := &config.AppConfig{
		JWT_SECRET: "test-secret",
	}

	middleware := JWTMiddleware(cfg)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware(handler)

	tests := []struct {
		name   string
		header string
	}{
		{"No Bearer prefix", "sometoken"},
		{"Wrong prefix", "Token sometoken"},
		{"Empty after Bearer", "Bearer "},
		{"Just Bearer", "Bearer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.header)

			w := httptest.NewRecorder()
			wrapped.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("Expected status 401, got %d", w.Code)
			}
		})
	}
}

func TestJWTMiddleware_ExpiredToken(t *testing.T) {
	secret := "test-secret"

	cfg := &config.AppConfig{
		JWT_SECRET: secret,
	}

	middleware := JWTMiddleware(cfg)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for expired token")
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware(handler)

	// Token expired 1 hour ago
	token := generateTestToken(secret, "user123", time.Now().Add(-time.Hour))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestJWTMiddleware_InvalidSignature(t *testing.T) {
	cfg := &config.AppConfig{
		JWT_SECRET: "correct-secret",
	}

	middleware := JWTMiddleware(cfg)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for invalid signature")
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware(handler)

	// Token signed with wrong secret
	token := generateTestToken("wrong-secret", "user123", time.Now().Add(time.Hour))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestJWTMiddleware_MalformedToken(t *testing.T) {
	cfg := &config.AppConfig{
		JWT_SECRET: "test-secret",
	}

	middleware := JWTMiddleware(cfg)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for malformed token")
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.jwt")

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
