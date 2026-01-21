package proxy

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"google.golang.org/api/idtoken"
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

	// Create authenticated HTTP client with identity token
	ctx := context.Background()
	client, err := idtoken.NewClient(ctx, cfg.TargetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticated client for %s: %w", cfg.TargetURL, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Use authenticated client instead of default transport
	proxy.Transport = client.Transport

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
			"fullURL", req.URL.String(),
		)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		slog.Info("Received response from backend",
			"service", cfg.Service,
			"status", resp.StatusCode,
			"contentLength", resp.ContentLength,
		)
		resp.Header.Del("Access-Control-Allow-Origin")
		resp.Header.Del("Access-Control-Allow-Methods")
		resp.Header.Del("Access-Control-Allow-Headers")
		resp.Header.Del("Access-Control-Allow-Credentials")
		resp.Header.Del("Access-Control-Max-Age")
		resp.Header.Del("Access-Control-Expose-Headers")

		return nil
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
	slog.Info("Creating USER SERVICE PROXY")
	return NewReverseProxy(ProxyConfig{
		TargetURL: targetURL,
		Service:   "user-service",
	})
}

// Proxy for session service
func NewSessionServiceProxy(targetURL string) (http.Handler, error) {
	slog.Info("Creating SESSION SERVICE PROXY")
	return NewReverseProxy(ProxyConfig{
		TargetURL: targetURL,
		Service:   "session-service",
	})
}

// Proxy for summary service
func NewSummaryServiceProxy(targetURL string) (http.Handler, error) {
	slog.Info("Creating SUMMARY SERVICE PROXY")
	return NewReverseProxy(ProxyConfig{
		TargetURL: targetURL,
		Service:   "summary-service",
	})
}

// Proxy for trend service
func NewTrendServiceProxy(targetURL string) (http.Handler, error) {
	slog.Info("Creating TREND SERVICE PROXY")
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
