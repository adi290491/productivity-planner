package proxy

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ProxyConfig holds configuration for reverse proxy
type ProxyConfig struct {
	TargetURL string
	Service   string
}

func NewReverseProxy(cfg ProxyConfig) (http.Handler, error) {
	targetURL, err := url.Parse(cfg.TargetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid target URL %s: %w", cfg.TargetURL, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Customize the Director to modify outbound requests
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		// Call original director
		originalDirector(req)

		// Add custom header
		// Extract userId from context
		if userID := req.Context().Value("userId"); userID != nil {
			req.Header.Set("X-User-ID", fmt.Sprintf("%v", userID))
		}

		slog.Debug("Proxying request",
			"path", req.URL.Path,
			"method", req.Method,
			"service", cfg.Service,
			"target", targetURL.String(),
		)
	}

	// Custom error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		slog.Error("Proxy error",
			"service", cfg.Service,
			"error", err,
			"method", r.Method,
			"path", r.URL.Path,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error":"service unavailable"}`))
	}
	return proxy, nil
}

// Proxy for user service
func NewUserServiceProxy(targetURL string) (http.Handler, error) {
	return NewReverseProxy(ProxyConfig{
		TargetURL: targetURL,
		Service:   "user-service",
	})
}

// Proxy for session service
func NewSessionServiceProxy(targetURL string) (http.Handler, error) {
	return NewReverseProxy(ProxyConfig{
		TargetURL: targetURL,
		Service:   "session-service",
	})
}

// Proxy for summary service
func NewSummaryServiceProxy(targetURL string) (http.Handler, error) {
	return NewReverseProxy(ProxyConfig{
		TargetURL: targetURL,
		Service:   "summary-service",
	})
}

// Proxy for trend service
func NewTrendServiceProxy(targetURL string) (http.Handler, error) {
	return NewReverseProxy(ProxyConfig{
		TargetURL: targetURL,
		Service:   "trend-service",
	})
}

// Custom ContextKey to avoid collisions
type ContextKey string

const (
	userIDKey ContextKey = "userId"
)

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID := ctx.Value("userId")
	if userID == nil {
		return "", false
	}

	switch v := userID.(type) {
	case string:
		return v, true
	default:
		return fmt.Sprintf("%v", v), true
	}
}
