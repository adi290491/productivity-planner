package proxy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// mockRoundTripper implements http.RoundTripper for testing
type mockRoundTripper struct {
	response *http.Response
	err      error
	mu       sync.Mutex
	requests []*http.Request
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	m.requests = append(m.requests, req)
	m.mu.Unlock()
	return m.response, m.err
}

func (m *mockRoundTripper) getRequests() []*http.Request {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Return a copy to avoid race conditions
	reqs := make([]*http.Request, len(m.requests))
	copy(reqs, m.requests)
	return reqs
}

func TestNewAuthenticatedTransport(t *testing.T) {
	tests := []struct {
		name        string
		audience    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid audience",
			audience: "https://example.com",
			wantErr:  false,
		},
		{
			name:        "empty audience",
			audience:    "",
			wantErr:     true,
			errContains: "audience cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport, err := NewAuthenticatedTransport(tt.audience)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errContains != "" && err != nil && err.Error() != tt.errContains {
					t.Errorf("error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if transport.audience != tt.audience {
				t.Errorf("audience = %v, want %v", transport.audience, tt.audience)
			}

			if transport.expiryBuffer != 5*time.Minute {
				t.Errorf("expiryBuffer = %v, want %v", transport.expiryBuffer, 5*time.Minute)
			}

			if transport.fetchTimeout != 10*time.Second {
				t.Errorf("fetchTimeout = %v, want %v", transport.fetchTimeout, 10*time.Second)
			}
		})
	}
}

func TestTransportOptions(t *testing.T) {
	customTransport := &mockRoundTripper{}

	transport, err := NewAuthenticatedTransport(
		"https://example.com",
		WithExpiryBuffer(10*time.Minute),
		WithFetchTimeout(30*time.Second),
		WithBaseTransport(customTransport),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.expiryBuffer != 10*time.Minute {
		t.Errorf("expiryBuffer = %v, want %v", transport.expiryBuffer, 10*time.Minute)
	}

	if transport.fetchTimeout != 30*time.Second {
		t.Errorf("fetchTimeout = %v, want %v", transport.fetchTimeout, 30*time.Second)
	}

	if transport.base != customTransport {
		t.Error("base transport not set correctly")
	}
}

func TestRoundTrip_Success(t *testing.T) {
	// Create mock response
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       http.NoBody,
		Header:     make(http.Header),
	}

	mockBase := &mockRoundTripper{
		response: mockResp,
	}

	// Create transport with mock token source
	transport, err := NewAuthenticatedTransport(
		"https://example.com",
		WithBaseTransport(mockBase),
	)
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	// Mock token fetcher
	expectedToken := "test-token-12345"
	transport.tokenFetcher = func(ctx context.Context, audience string) (string, time.Time, error) {
		return expectedToken, time.Now().Add(1 * time.Hour), nil
	}

	// Create test request
	req := httptest.NewRequest("GET", "https://example.com/test", nil)

	// Execute RoundTrip
	resp, err := transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	// Verify authorization header was added
	requests := mockBase.getRequests()
	if len(requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(requests))
	}

	authHeader := requests[0].Header.Get("Authorization")
	expectedAuth := "Bearer " + expectedToken
	if authHeader != expectedAuth {
		t.Errorf("Authorization = %v, want %v", authHeader, expectedAuth)
	}
}

func TestRoundTrip_TokenError(t *testing.T) {
	mockBase := &mockRoundTripper{
		response: &http.Response{StatusCode: http.StatusOK},
	}

	transport, err := NewAuthenticatedTransport(
		"https://example.com",
		WithBaseTransport(mockBase),
	)
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	// Mock token fetcher that returns error
	expectedErr := errors.New("token fetch failed")
	transport.tokenFetcher = func(ctx context.Context, audience string) (string, time.Time, error) {
		return "", time.Time{}, expectedErr
	}

	req := httptest.NewRequest("GET", "https://example.com/test", nil)

	_, err = transport.RoundTrip(req)
	if err == nil {
		t.Error("expected error, got nil")
	}

	// Verify no request was made to backend
	if len(mockBase.requests) != 0 {
		t.Errorf("expected 0 requests to backend, got %d", len(mockBase.requests))
	}
}

func TestRoundTrip_RequestError(t *testing.T) {
	expectedErr := errors.New("network error")
	mockBase := &mockRoundTripper{
		err: expectedErr,
	}

	transport, err := NewAuthenticatedTransport(
		"https://example.com",
		WithBaseTransport(mockBase),
	)
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	// Mock successful token fetcher
	transport.tokenFetcher = func(ctx context.Context, audience string) (string, time.Time, error) {
		return "test-token", time.Now().Add(1 * time.Hour), nil
	}

	req := httptest.NewRequest("GET", "https://example.com/test", nil)

	_, err = transport.RoundTrip(req)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("error = %v, want %v", err, expectedErr)
	}

	// The request IS made to backend, it just fails
	// Verify the request was attempted (with auth header)
	requests := mockBase.getRequests()
	if len(requests) != 1 {
		t.Errorf("expected 1 request attempt to backend, got %d", len(requests))
	}

	// Verify auth header was added even though request failed
	if len(requests) > 0 {
		auth := requests[0].Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Authorization header even on failed request, got %s", auth)
		}
	}
}

func TestTokenCaching(t *testing.T) {
	mockBase := &mockRoundTripper{
		response: &http.Response{StatusCode: http.StatusOK, Body: http.NoBody},
	}

	transport, err := NewAuthenticatedTransport(
		"https://example.com",
		WithBaseTransport(mockBase),
		WithExpiryBuffer(1*time.Minute),
	)
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	// Track how many times token is fetched (thread-safe)
	var fetchCount int32
	transport.tokenFetcher = func(ctx context.Context, audience string) (string, time.Time, error) {
		atomic.AddInt32(&fetchCount, 1)
		return "test-token", time.Now().Add(10 * time.Minute), nil
	}

	// Make first request
	req1 := httptest.NewRequest("GET", "https://example.com/test", nil)
	_, err = transport.RoundTrip(req1)
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	if atomic.LoadInt32(&fetchCount) != 1 {
		t.Errorf("fetchCount after first request = %d, want 1", atomic.LoadInt32(&fetchCount))
	}

	// Make second request - should use cached token
	req2 := httptest.NewRequest("GET", "https://example.com/test2", nil)
	_, err = transport.RoundTrip(req2)
	if err != nil {
		t.Fatalf("second request failed: %v", err)
	}

	if atomic.LoadInt32(&fetchCount) != 1 {
		t.Errorf("fetchCount after second request = %d, want 1 (should use cache)", atomic.LoadInt32(&fetchCount))
	}

	// Verify both requests used the same token
	requests := mockBase.getRequests()
	if len(requests) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(requests))
	}

	token1 := requests[0].Header.Get("Authorization")
	token2 := requests[1].Header.Get("Authorization")

	if token1 != token2 {
		t.Errorf("tokens differ: %v vs %v", token1, token2)
	}
}

func TestTokenRefresh(t *testing.T) {
	mockBase := &mockRoundTripper{
		response: &http.Response{StatusCode: http.StatusOK, Body: http.NoBody},
	}

	transport, err := NewAuthenticatedTransport(
		"https://example.com",
		WithBaseTransport(mockBase),
		WithExpiryBuffer(1*time.Minute),
	)
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	// Track tokens returned
	tokenIndex := 0
	tokens := []string{"token-1", "token-2"}

	transport.tokenFetcher = func(ctx context.Context, audience string) (string, time.Time, error) {
		token := tokens[tokenIndex]
		tokenIndex++
		return token, time.Now().Add(2 * time.Second), nil
	}

	// First request
	req1 := httptest.NewRequest("GET", "https://example.com/test", nil)
	_, err = transport.RoundTrip(req1)
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	requests := mockBase.getRequests()
	if len(requests) < 1 {
		t.Fatal("expected at least 1 request")
	}

	auth1 := requests[0].Header.Get("Authorization")
	expectedAuth1 := "Bearer token-1"
	if auth1 != expectedAuth1 {
		t.Errorf("first auth = %v, want %v", auth1, expectedAuth1)
	}

	// Wait for token to expire (expiry - buffer = 2s - 1m < 0, so it's already expired)
	time.Sleep(100 * time.Millisecond)

	// Second request should get new token
	req2 := httptest.NewRequest("GET", "https://example.com/test2", nil)
	_, err = transport.RoundTrip(req2)
	if err != nil {
		t.Fatalf("second request failed: %v", err)
	}

	requests = mockBase.getRequests()
	if len(requests) < 2 {
		t.Fatal("expected at least 2 requests")
	}

	auth2 := requests[1].Header.Get("Authorization")
	expectedAuth2 := "Bearer token-2"
	if auth2 != expectedAuth2 {
		t.Errorf("second auth = %v, want %v", auth2, expectedAuth2)
	}
}

func TestClearCache(t *testing.T) {
	transport, err := NewAuthenticatedTransport("https://example.com")
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	// Set up token fetcher
	transport.tokenFetcher = func(ctx context.Context, audience string) (string, time.Time, error) {
		return "test-token", time.Now().Add(1 * time.Hour), nil
	}

	// Fetch token
	_, err = transport.getToken(context.Background())
	if err != nil {
		t.Fatalf("failed to get token: %v", err)
	}

	if !transport.HasValidToken() {
		t.Error("should have valid token after fetch")
	}

	// Clear cache
	transport.ClearCache()

	if transport.HasValidToken() {
		t.Error("should not have valid token after clear")
	}

	if !transport.GetCachedTokenExpiry().IsZero() {
		t.Error("expiry should be zero after clear")
	}
}

func TestConcurrentAccess(t *testing.T) {
	mockBase := &mockRoundTripper{
		response: &http.Response{StatusCode: http.StatusOK, Body: http.NoBody},
	}

	transport, err := NewAuthenticatedTransport(
		"https://example.com",
		WithBaseTransport(mockBase),
	)
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	var fetchCount int32

	transport.tokenFetcher = func(ctx context.Context, audience string) (string, time.Time, error) {
		atomic.AddInt32(&fetchCount, 1)

		// Simulate slow token fetch
		time.Sleep(100 * time.Millisecond)

		return "test-token", time.Now().Add(1 * time.Hour), nil
	}

	// Launch 10 concurrent requests
	const numRequests = 10
	errChan := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(idx int) {
			req := httptest.NewRequest("GET", "https://example.com/test", nil)
			_, err := transport.RoundTrip(req)
			errChan <- err
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		if err := <-errChan; err != nil {
			t.Errorf("request %d failed: %v", i, err)
		}
	}

	// Token should only be fetched once despite concurrent requests
	count := atomic.LoadInt32(&fetchCount)
	if count != 1 {
		t.Errorf("fetchCount = %d, want 1 (should only fetch once for concurrent requests)", count)
	}
}

func TestEmptyToken(t *testing.T) {
	transport, err := NewAuthenticatedTransport("https://example.com")
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	// Mock token fetcher that returns empty token
	transport.tokenFetcher = func(ctx context.Context, audience string) (string, time.Time, error) {
		return "", time.Now().Add(1 * time.Hour), fmt.Errorf("received empty token")
	}

	_, err = transport.getToken(context.Background())
	if err == nil {
		t.Error("expected error for empty token, got nil")
	}

	if err.Error() != "received empty token" {
		t.Errorf("error = %v, want 'received empty token'", err)
	}
}

func TestContextCancellation(t *testing.T) {
	transport, err := NewAuthenticatedTransport(
		"https://example.com",
		WithFetchTimeout(5*time.Second),
	)
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	// Mock token fetcher that takes too long
	transport.tokenFetcher = func(ctx context.Context, audience string) (string, time.Time, error) {
		select {
		case <-ctx.Done():
			return "", time.Time{}, ctx.Err()
		case <-time.After(10 * time.Second):
			return "test-token", time.Now().Add(1 * time.Hour), nil
		}
	}

	// Create context that cancels immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = transport.getToken(ctx)
	if err == nil {
		t.Error("expected error from cancelled context, got nil")
	}
}
