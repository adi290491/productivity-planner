package router

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/adi290491/productivity-planner/gateway/config"
	"github.com/adi290491/productivity-planner/gateway/middleware"
	"github.com/adi290491/productivity-planner/gateway/proxy"
	"github.com/adi290491/productivity-planner/gateway/ratelimiter"
)

// NewRouter creates and configures an HTTP router with all routes and middleware
func NewRouter(cfg *config.AppConfig, rateLimiterMgr *ratelimiter.RateLimiterManager, rateLimitCfg *ratelimiter.Config) http.Handler {
	mux := http.NewServeMux()

	// Middleware setup
	corsMiddleware := middleware.CorsMiddleware()
	rateLimitMiddleware := ratelimiter.Middleware(rateLimiterMgr, rateLimitCfg)
	jwtMiddleware := middleware.JWTMiddleware(cfg)

	// Service URLs from env
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	sessionServiceURL := os.Getenv("SESSION_SERVICE_URL")
	summaryServiceURL := os.Getenv("SUMMARY_SERVICE_URL")
	trendServiceURL := os.Getenv("TREND_SERVICE_URL")

	userHandler := createProxyHandler("user-service", userServiceURL, proxy.NewUserServiceProxy)
	sessionHandler := createProxyHandler("session-service", sessionServiceURL, proxy.NewSessionServiceProxy)
	summaryHandler := createProxyHandler("summary-service", summaryServiceURL, proxy.NewSummaryServiceProxy)
	trendHandler := createProxyHandler("trend-service", trendServiceURL, proxy.NewTrendServiceProxy)

	// Public routes
	// user-service
	mux.Handle("/users/", middleware.Chain(
		userHandler,
		corsMiddleware,
		rateLimitMiddleware,
	))

	// JWT protected routes
	// session-service
	mux.Handle("/sessions/", middleware.Chain(
		sessionHandler,
		corsMiddleware,
		rateLimitMiddleware,
		jwtMiddleware,
	))

	// JWT protected routes
	// summary service
	mux.Handle("/summary/", middleware.Chain(
		summaryHandler,
		corsMiddleware,
		rateLimitMiddleware,
		jwtMiddleware,
	))

	// JWT protected routes
	// trend-service
	mux.Handle("/trend/", middleware.Chain(
		trendHandler,
		corsMiddleware,
		rateLimitMiddleware,
		jwtMiddleware,
	))

	return mux

}

func createProxyHandler(serviceName, serviceURL string, proxyFactory func(string) (http.Handler, error)) http.Handler {
	if serviceURL == "" {
		slog.Warn("Service URL not configured, using placeholder",
			"service", serviceName,
		)
		return createPlaceholderHandler(serviceName, "URL not configured")
	}

	handler, err := proxyFactory(serviceURL)
	if err != nil {
		slog.Error("Failed to create proxy, using placeholder",
			"service", serviceName,
			"error", err,
		)
		return createPlaceholderHandler(serviceName, serviceURL)
	}

	slog.Info("Proxy created",
		"service", serviceName,
		"target", serviceURL,
	)

	return handler
}

func createPlaceholderHandler(serviceName, serviceURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"` + serviceName + ` - placeholder", "url":"` + serviceURL + `"}`))
	})
}
