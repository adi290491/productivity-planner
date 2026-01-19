package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func CorsMiddleware() func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Define allowed origins (fallback for local development)
			allowedOrigins := []string{
				"https://systemic-productivity-planner.web.app",
				"https://systemic-productivity-planner.firebaseapp.com",
				"http://localhost:4200",
				"http://localhost:5173",
				"http://localhost:3000",
			}

			// Add environment-based origins (production)
			if frontendOrigin := os.Getenv("FRONTEND_ORIGIN"); frontendOrigin != "" {
				allowedOrigins = append([]string{frontendOrigin}, allowedOrigins...)
			}

			if frontendOrigin2 := os.Getenv("FRONTEND_ORIGIN_2"); frontendOrigin2 != "" {
				allowedOrigins = append([]string{frontendOrigin2}, allowedOrigins...)
			}

			// Check if origin is in allowed list
			originAllowed := false
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					originAllowed = true
					break
				}
			}

			// Always set other CORS headers for valid origins
			if originAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-USER-ID")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			if r.Method == http.MethodOptions {
				if originAllowed {
					w.WriteHeader(http.StatusNoContent)
				} else {
					slog.Warn("CORS preflight rejected", "origin", origin, "path", r.URL.Path)
					w.WriteHeader(http.StatusForbidden)
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// extractOriginDomain extracts the domain from an origin
func extractOriginDomain(origin string) string {
	// Remove protocol
	origin = strings.TrimPrefix(origin, "https://")
	origin = strings.TrimPrefix(origin, "http://")

	// Remove port if present
	if colonIdx := strings.Index(origin, ":"); colonIdx != -1 {
		origin = origin[:colonIdx]
	}
	return origin
}
