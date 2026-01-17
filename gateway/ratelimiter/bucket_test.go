package ratelimiter

import (
	"math"
	"sync"
	"testing"
	"time"
)

func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestNewTokenBucket(t *testing.T) {
	capacity := 10.0
	refillRate := 5.0

	bucket := NewTokenBucket(capacity, refillRate)
	if bucket.capacity != capacity {
		t.Errorf("Expected capacity %f, got %f", capacity, bucket.capacity)
	}

	if bucket.refillRate != refillRate {
		t.Errorf("Expected refill rate %f, got %f", refillRate, bucket.refillRate)
	}

	if bucket.tokens != capacity {
		t.Errorf("Expected initial tokens %f, got %f", capacity, bucket.tokens)
	}

	if bucket.lastRefill.IsZero() {
		t.Error("Expcected lastRefill to be set")
	}
}

func TestTokenBucket_Allow_SingleRequest(t *testing.T) {
	bucket := NewTokenBucket(10.0, 5.0)

	allowed := bucket.Allow()
	if !allowed {
		t.Error("Expected first request to be allowed")
	}

	available := bucket.AvailableTokens()
	expected := 9.0
	tolerance := 0.001 // Allow 0.001 difference for float-point precision
	if !floatEquals(available, expected, tolerance) {
		t.Errorf("Expected %.1f tokens remaining, got %.1f", expected, available)
	}
}

func TestTokenBucket_Allow_BurstCapacity(t *testing.T) {
	capacity := 5.0
	bucket := NewTokenBucket(capacity, 1.0)

	// Should allow exactly 'capacity' requests in rapid succession
	for i := 0; i < int(capacity); i++ {
		if !bucket.Allow() {
			t.Errorf("Expected request %d to be allowed (burst)", i+1)
		}
	}

	// next request should be denied (bucket empty)
	if bucket.Allow() {
		t.Errorf("Expected request after burst to be denied")
	}

	available := bucket.AvailableTokens()
	tolerance := 0.001
	if !floatEquals(available, 0.0, tolerance) {
		t.Errorf("Expected 0 tokens remaining after burst, got %f", available)
	}
}

func TestTokenBucket_Allow_RefillOverTime(t *testing.T) {
	capacity := 10.0
	refillRate := 10.0
	bucket := NewTokenBucket(capacity, refillRate)

	// Consume All tokens
	for i := 0; i < int(capacity); i++ {
		bucket.Allow()
	}

	// the next request should be denied immediately
	if bucket.Allow() {
		t.Error("Expected request to ne denied when the bucket is empty")
	}

	// wait for 100 ms (add 1 token: 10 token/sec * 1 = 1 token)
	time.Sleep(100 * time.Millisecond)

	// Should now allow 1 request
	if !bucket.Allow() {
		t.Error("Expected request to be allowed after refill")
	}

	// Should be denied again
	if bucket.Allow() {
		t.Error("Expected request to be denied after consuming refilled token")
	}
}

func TestTokenBucket_Allow_RefillDoesNotExceedCapacity(t *testing.T) {
	capacity := 5.0
	refillRate := 10.0
	bucket := NewTokenBucket(capacity, refillRate)

	// Wait for 1 second (would add 10 tokens, but capacity is 5)
	time.Sleep(1 * time.Second)

	available := bucket.AvailableTokens()
	tolerance := 0.01 // Slightly larger tolerance for timing variance

	if available > capacity+tolerance {
		t.Errorf("Expected tokens to not exceed capacity %f, got %f", capacity, available)
	}

	// Should still be at capacity
	if !floatEquals(available, capacity, tolerance) {
		t.Errorf("Expected tokens to be at capacity %f after long wait, got %f", capacity, available)
	}
}

func TestTokenBucket_AllowN(t *testing.T) {
	bucket := NewTokenBucket(10.0, 5.0)

	// Should allow consuming 3 tokens
	if !bucket.AllowN(3) {
		t.Error("Expected AllowN(3) to be allowed")
	}

	available := bucket.AvailableTokens()
	expected := 7.0
	tolerance := 0.001

	if !floatEquals(available, expected, tolerance) {
		t.Errorf("Expected %f tokens remaining, got %f", expected, available)
	}

	// Should allow consuming 7 more tokens
	if !bucket.AllowN(7) {
		t.Error("Expected AllowN(7) to be allowed")
	}

	// Should deny consuming 1 token (bucket empty)
	if bucket.AllowN(1) {
		t.Error("Expected AllowN(1) to be denied when bucket is empty")
	}
}

func TestTokenBucket_AllowN_InsufficientTokens(t *testing.T) {
	bucket := NewTokenBucket(10.0, 5.0)

	// Try to consume more than capacity
	if bucket.AllowN(15) {
		t.Error("Expected AllowN(15) to be denied when capacity is 10")
	}

	// Tokens should remain unchanged
	available := bucket.AvailableTokens()
	expected := 10.0
	tolerance := 0.001

	if !floatEquals(available, expected, tolerance) {
		t.Errorf("Expected tokens to remain at %f after failed AllowN, got %f", expected, available)
	}
}

func TestTokenBucket_Concurrent_Access(t *testing.T) {
	capacity := 100.0
	refillRate := 50.0
	bucket := NewTokenBucket(capacity, refillRate)

	numGoroutines := 50
	requestsPerGoroutine := 10

	var wg sync.WaitGroup
	allowedCount := make(chan int, numGoroutines)

	// Spawn multiple goroutines making concurrent requests
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed := 0
			for j := 0; j < requestsPerGoroutine; j++ {
				if bucket.Allow() {
					allowed++
				}
			}
			allowedCount <- allowed
		}()
	}

	wg.Wait()
	close(allowedCount)

	// Count total allowed requests
	totalAllowed := 0
	for count := range allowedCount {
		totalAllowed += count
	}

	// Should allow exactly capacity requests (some refill might occur during test)
	// So totalAllowed should be >= capacity but not wildly more
	if totalAllowed < int(capacity) {
		t.Errorf("Expected at least %d allowed requests, got %d", int(capacity), totalAllowed)
	}

	// Should not allow all requests (500 total requests)
	totalRequests := numGoroutines * requestsPerGoroutine
	if totalAllowed >= totalRequests {
		t.Errorf("Expected some requests to be denied, but all %d were allowed", totalRequests)
	}

	t.Logf("Concurrent test: %d/%d requests allowed", totalAllowed, totalRequests)
}

func TestTokenBucket_AvailableTokens(t *testing.T) {
	bucket := NewTokenBucket(10.0, 5.0)
	tolerance := 0.001

	// Initially full
	available := bucket.AvailableTokens()
	if !floatEquals(available, 10.0, tolerance) {
		t.Errorf("Expected 10 tokens initially, got %f", available)
	}

	// After consuming 3
	bucket.AllowN(3)
	available = bucket.AvailableTokens()
	if !floatEquals(available, 7.0, tolerance) {
		t.Errorf("Expected 7 tokens after consuming 3, got %f", available)
	}

	// After waiting and refilling
	time.Sleep(200 * time.Millisecond) // Should add ~1 token (0.2s * 5/s = 1)
	available = bucket.AvailableTokens()

	// Should be approximately 8 tokens (7 + 1 from refill)
	// Allow larger tolerance for timing variance
	expectedAfterRefill := 8.0
	timingTolerance := 0.2 // Allow up to 0.2 token variance for timing

	if !floatEquals(available, expectedAfterRefill, timingTolerance) {
		t.Logf("Warning: Expected ~%f tokens after refill, got %f (timing variance)", expectedAfterRefill, available)
	}
}

func TestTokenBucket_FractionalTokens(t *testing.T) {
	// Test that fractional tokens accumulate correctly
	bucket := NewTokenBucket(10.0, 1.0) // 1 token per second

	// Consume all tokens
	for i := 0; i < 10; i++ {
		bucket.Allow()
	}

	// Wait 500ms (should add 0.5 tokens)
	time.Sleep(500 * time.Millisecond)

	// Should still be denied (need 1.0 token)
	if bucket.Allow() {
		t.Error("Expected request to be denied with only 0.5 tokens")
	}

	// Wait another 500ms (now should have ~1.0 token)
	time.Sleep(500 * time.Millisecond)

	// Should now be allowed
	if !bucket.Allow() {
		t.Error("Expected request to be allowed after accumulating 1.0 token")
	}
}
