package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func init() {
	// Suppress logs during tests
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestCorsMiddleware_AllowedOrigin(t *testing.T) {
	middleware := CorsMiddleware()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	wrapped := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://systemic-productivity-planner.web.app")

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check CORS headers
	tests := map[string]string{
		"Access-Control-Allow-Origin":      "https://systemic-productivity-planner.web.app",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization, X-USER-ID",
		"Access-Control-Allow-Credentials": "true",
	}

	for header, expected := range tests {
		if got := w.Header().Get(header); got != expected {
			t.Errorf("Header %s: expected %s, got %s", header, expected, got)
		}
	}
}

func TestCorsMiddleware_DisallowedOrigin(t *testing.T) {
	middleware := CorsMiddleware()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://malicious-site.com")

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	// Should still process the request but without CORS headers
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Should NOT have Access-Control-Allow-Origin for disallowed origin
	if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "" {
		t.Errorf("Expected no Access-Control-Allow-Origin for disallowed origin, got %s", origin)
	}
}

func TestCorsMiddleware_PreflightAllowed(t *testing.T) {
	middleware := CorsMiddleware()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for OPTIONS preflight")
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware(handler)

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Check CORS headers are present
	if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:5173" {
		t.Errorf("Expected Access-Control-Allow-Origin: http://localhost:5173, got %s", origin)
	}
}

func TestCorsMiddleware_PreflightDisallowed(t *testing.T) {
	middleware := CorsMiddleware()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for OPTIONS preflight")
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware(handler)

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://evil.com")

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestCorsMiddleware_EnvironmentOrigin(t *testing.T) {
	// Set environment variable
	os.Setenv("FRONTEND_ORIGIN", "https://custom-domain.com")
	defer os.Unsetenv("FRONTEND_ORIGIN")

	middleware := CorsMiddleware()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://custom-domain.com")

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "https://custom-domain.com" {
		t.Errorf("Expected custom origin to be allowed, got %s", origin)
	}
}

func TestCorsMiddleware_NoOriginHeader(t *testing.T) {
	middleware := CorsMiddleware()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	// No Origin header set

	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	// Should still work, just no CORS headers
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
