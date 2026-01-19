package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/adi290491/productivity-planner/gateway/config"
)

// JWTMiddleware creates a JWT authentication middleware
func JWTMiddleware(cfg *config.AppConfig) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")

			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				slog.Warn("Missing or invalid Authorization header",
					"path", r.URL.Path,
					"method", r.Method,
					"ip", r.RemoteAddr,
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Missing or invalid auth token",
				})
				return
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			if tokenStr == "" {
				slog.Warn("Empty auth token",
					"path", r.URL.Path,
					"method", r.Method,
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Empty auth token",
				})
				return
			}

			jwtUtil := JWTUtil{
				Secret:      cfg.JWT_SECRET,
				tokenString: tokenStr,
			}

			userId, err := jwtUtil.ValidateToken()
			if err != nil {
				slog.Warn("Invalid token",
					"path", r.URL.Path,
					"method", r.Method,
					"error", err,
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Invalid token",
				})
				return
			}

			if userId == nil {
				slog.Warn("Invalid user ID in token",
					"path", r.URL.Path,
					"method", r.Method,
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Invalid user ID in token",
				})
				return
			}

			slog.Debug("JWT validation successful",
				"path", r.URL.Path,
				"method", r.Method,
				"userId", userId,
			)

			// Add userId to request context
			ctx := context.WithValue(r.Context(), "userId", userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
