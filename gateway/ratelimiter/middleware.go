package ratelimiter

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func Middleware(manager *RateLimiterManager, config *Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if ratelimiting is disabled
			if !config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract client IP
			clientIP := extractClientIP(r)
			key := "ip:" + clientIP

			// Log BEFORE checking rate limit
			currentTokens := manager.AvailableTokens(key)
			slog.Info("Rate limit check START",
				"ip", clientIP,
				"path", r.URL.Path,
				"method", r.Method,
				"tokensBefore", currentTokens,
			)

			allowed, availableTokens := manager.AllowWithRemaining(key)

			slog.Info("Rate limit check END",
				"ip", clientIP,
				"path", r.URL.Path,
				"method", r.Method,
				"allowed", allowed,
				"tokensAfter", availableTokens,
			)

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", config.Capacity))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%.0f", availableTokens))

			if !allowed {
				//Rate limit exceeded
				slog.Warn("Rate limit exceeded",
					"ip", clientIP,
					"path", r.URL.Path,
					"method", r.Method,
					"remaining", int(availableTokens),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)

				response := map[string]interface{}{
					"error":   "rate limit exceeded",
					"message": "too many requests, please try again later",
				}

				json.NewEncoder(w).Encode(response)
				return
			}

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// extractClientIP extracts the client IP address from the request
// Checks X-Forwarded-For header first (for requests behind proxies/load balancers)
// Falls back to RemoteAddr
func extractClientIP(r *http.Request) string {

	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs
		// We need the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header (alternative header)
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fallback to RemoteAddr
	// RemoteAddr is in format "IP:port", we need just the IP
	ip := r.RemoteAddr
	if colonIdx := strings.LastIndex(ip, ":"); colonIdx != -1 {
		ip = ip[:colonIdx]
	}

	return ip

}
