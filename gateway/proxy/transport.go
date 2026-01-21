package proxy

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/api/idtoken"
)

type AuthenticatedTransport struct {
	base     http.RoundTripper
	audience string

	mu          sync.RWMutex
	cachedToken string
	tokenExpiry time.Time

	// Configuration
	expiryBuffer time.Duration // How much earlier to refresh token
	fetchTimeout time.Duration // Timeout for fetching tokens

	tokenFetcher func(ctx context.Context, audience string) (string, time.Time, error)
}

// TransportOption configures the AuthenticatedTransport
type TransportOption func(*AuthenticatedTransport)

// WithExpiryBuffer sets how much before expiry to refresh the token
func WithExpiryBuffer(d time.Duration) TransportOption {
	return func(t *AuthenticatedTransport) {
		t.expiryBuffer = d
	}
}

// WithFetchTimeout sets the timeout for fetching tokens
func WithFetchTimeout(d time.Duration) TransportOption {
	return func(t *AuthenticatedTransport) {
		t.fetchTimeout = d
	}
}

// WithBaseTransport sets the underlying transport
func WithBaseTransport(rt http.RoundTripper) TransportOption {
	return func(t *AuthenticatedTransport) {
		t.base = rt
	}
}

// defaultTokenFetcher fetches a token using the Google Cloud idtoken library
func defaultTokenFetcher(ctx context.Context, audience string) (string, time.Time, error) {
	tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to create token source: %w", err)
	}

	token, err := tokenSource.Token()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get token: %w", err)
	}

	if token.AccessToken == "" {
		return "", time.Time{}, fmt.Errorf("received empty token")
	}

	return token.AccessToken, token.Expiry, nil
}

// NewAuthenticatedTransport creates a new transport that adds identity tokens to requests
func NewAuthenticatedTransport(audience string, opts ...TransportOption) (*AuthenticatedTransport, error) {
	if audience == "" {
		return nil, fmt.Errorf("audience cannot be empty")
	}

	baseTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	t := &AuthenticatedTransport{
		base:         baseTransport,
		audience:     audience,
		expiryBuffer: 5 * time.Minute,
		fetchTimeout: 10 * time.Second,
		tokenFetcher: defaultTokenFetcher,
	}

	for _, opt := range opts {
		opt(t)
	}

	return t, nil
}

// RoundTrip implements http.RoundTripper
func (t *AuthenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	reqClone := req.Clone(req.Context())

	// Get or refresh the token
	token, err := t.getToken(req.Context())
	if err != nil {
		slog.Error("Failed to get identity token",
			"audience", t.audience,
			"error", err,
			"path", req.URL.Path,
		)
		return nil, fmt.Errorf("failed to get identity token: %w", err)
	}

	reqClone.Header.Set("Authorization", "Bearer "+token)

	slog.Debug("Added identity token to request",
		"audience", t.audience,
		"path", req.URL.Path,
		"method", req.Method,
	)

	// Make the request with the base transport
	resp, err := t.base.RoundTrip(reqClone)
	if err != nil {
		slog.Error("Request failed",
			"audience", t.audience,
			"error", err,
			"path", req.URL.Path,
		)
		return nil, err
	}

	slog.Debug("Received response",
		"audience", t.audience,
		"status", resp.StatusCode,
		"path", req.URL.Path,
	)

	return resp, nil
}

// getToken returns a cached token or fetches a new one if expired
func (t *AuthenticatedTransport) getToken(ctx context.Context) (string, error) {

	t.mu.RLock()
	if t.cachedToken != "" && time.Now().Before(t.tokenExpiry) {
		token := t.cachedToken
		t.mu.RUnlock()
		slog.Debug("Using cached identity token",
			"audience", t.audience,
			"expiresIn", time.Until(t.tokenExpiry).Round(time.Second),
		)
		return token, nil
	}
	t.mu.RUnlock()

	return t.refreshToken(ctx)
}

// refreshToken fetches a new token and caches it
func (t *AuthenticatedTransport) refreshToken(ctx context.Context) (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Double-check in case another goroutine just fetched it
	if t.cachedToken != "" && time.Now().Before(t.tokenExpiry) {
		slog.Debug("Token was refreshed by another goroutine",
			"audience", t.audience,
		)
		return t.cachedToken, nil
	}

	slog.Info("Fetching new identity token",
		"audience", t.audience,
	)

	// Create context with timeout
	fetchCtx, cancel := context.WithTimeout(ctx, t.fetchTimeout)
	defer cancel()

	// Fetch new token
	accessToken, expiry, err := t.tokenFetcher(fetchCtx, t.audience)
	if err != nil {
		return "", err
	}

	// Cache the token (expire buffer time before actual expiry for safety)
	t.cachedToken = accessToken
	t.tokenExpiry = expiry.Add(-t.expiryBuffer)

	slog.Info("Successfully fetched new identity token",
		"audience", t.audience,
		"expiresAt", expiry.Format(time.RFC3339),
		"expiresIn", time.Until(t.tokenExpiry).Round(time.Second),
	)

	return t.cachedToken, nil
}

// ClearCache clears the cached token (useful for testing or forcing refresh)
func (t *AuthenticatedTransport) ClearCache() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.cachedToken = ""
	t.tokenExpiry = time.Time{}

	slog.Debug("Cleared token cache",
		"audience", t.audience,
	)
}

// GetCachedTokenExpiry returns the expiry time of the cached token (for testing/monitoring)
func (t *AuthenticatedTransport) GetCachedTokenExpiry() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.tokenExpiry
}

// HasValidToken returns true if there's a valid cached token
func (t *AuthenticatedTransport) HasValidToken() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.cachedToken != "" && time.Now().Before(t.tokenExpiry)
}
