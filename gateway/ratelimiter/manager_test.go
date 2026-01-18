package ratelimiter

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewRateLimiterManager(t *testing.T) {
	capacity := 10.0
	refillRate := 5.0
	cleanupInterval := 1 * time.Minute
	ttl := 30 * time.Minute

	manager := NewRateLimiterManager(capacity, refillRate, cleanupInterval, ttl)
	defer manager.Shutdown()

	if manager.capacity != capacity {
		t.Errorf("Expected capacity %f, got %f", capacity, manager.capacity)
	}

	if manager.refillRate != refillRate {
		t.Errorf("Expected refillRate %f, got %f", refillRate, manager.refillRate)
	}

	if manager.cleanupInterval != cleanupInterval {
		t.Errorf("Expected cleanupInterval %v, got %v", cleanupInterval, manager.cleanupInterval)
	}

	if manager.ttl != ttl {
		t.Errorf("Expected ttl %v, got %v", ttl, manager.ttl)
	}
}

func TestRateLimiterManager_GetLimiter(t *testing.T) {
	manager := NewRateLimiterManager(10.0, 5.0, 1*time.Minute, 30*time.Minute)
	defer manager.Shutdown()

	key := "192.168.1.1"

	// Get limiter for first time
	bucket1 := manager.GetLimiter(key)
	if bucket1 == nil {
		t.Fatal("Expected non-nil bucket")
	}

	// Get same limiter again - should be the same instance
	bucket2 := manager.GetLimiter(key)
	if bucket1 != bucket2 {
		t.Error("Expected same bucket instance for same key")
	}

	// Different key should get different bucket
	key2 := "192.168.1.2"
	bucket3 := manager.GetLimiter(key2)
	if bucket1 == bucket3 {
		t.Error("Expected different bucket for different key")
	}
}

func TestRateLimiterManager_Allow(t *testing.T) {
	manager := NewRateLimiterManager(5.0, 1.0, 1*time.Minute, 30*time.Minute)
	defer manager.Shutdown()

	key := "user:123"

	// Should allow up to capacity
	for i := 0; i < 5; i++ {
		if !manager.Allow(key) {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	// Next should be denied
	if manager.Allow(key) {
		t.Error("Expected request to be denied after capacity exhausted")
	}

	// Different key should have its own bucket
	key2 := "user:456"
	if !manager.Allow(key2) {
		t.Error("Expected request for different key to be allowed")
	}
}

func TestRateLimiterManager_AllowN(t *testing.T) {
	manager := NewRateLimiterManager(10.0, 5.0, 1*time.Minute, 30*time.Minute)
	defer manager.Shutdown()

	key := "api-key-123"

	// Should allow consuming 3 tokens
	if !manager.AllowN(key, 3) {
		t.Error("Expected AllowN(3) to be allowed")
	}

	// Should allow consuming 7 more tokens
	if !manager.AllowN(key, 7) {
		t.Error("Expected AllowN(7) to be allowed")
	}

	// Should deny consuming 1 token (bucket empty)
	if manager.AllowN(key, 1) {
		t.Error("Expected AllowN(1) to be denied when bucket is empty")
	}
}

func TestRateLimiterManager_AvailableTokens(t *testing.T) {
	manager := NewRateLimiterManager(10.0, 5.0, 1*time.Minute, 30*time.Minute)
	defer manager.Shutdown()

	key := "test-key"
	tolerance := 0.001

	// Initially should have full capacity
	available := manager.AvailableTokens(key)
	if !floatEquals(available, 10.0, tolerance) {
		t.Errorf("Expected 10 tokens initially, got %f", available)
	}

	// After consuming 3
	manager.AllowN(key, 3)
	available = manager.AvailableTokens(key)
	if !floatEquals(available, 7.0, tolerance) {
		t.Errorf("Expected 7 tokens after consuming 3, got %f", available)
	}
}

func TestRateLimiterManager_Cleanup(t *testing.T) {
	// Use short intervals for testing
	cleanupInterval := 100 * time.Millisecond
	ttl := 200 * time.Millisecond

	manager := NewRateLimiterManager(10.0, 5.0, cleanupInterval, ttl)
	defer manager.Shutdown()

	// Create several buckets
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key-%d", i)
		manager.Allow(key)
	}

	// Should have 5 buckets
	if count := manager.Count(); count != 5 {
		t.Errorf("Expected 5 buckets, got %d", count)
	}

	// Wait for TTL to expire
	time.Sleep(ttl + 50*time.Millisecond)

	// Wait for cleanup to run (cleanup interval + buffer)
	time.Sleep(cleanupInterval + 100*time.Millisecond)

	// All buckets should be cleaned up
	if count := manager.Count(); count != 0 {
		t.Errorf("Expected 0 buckets after cleanup, got %d", count)
	}
}

func TestRateLimiterManager_CleanupKeepsActive(t *testing.T) {
	cleanupInterval := 100 * time.Millisecond
	ttl := 200 * time.Millisecond

	manager := NewRateLimiterManager(10.0, 5.0, cleanupInterval, ttl)
	defer manager.Shutdown()

	activeKey := "active-key"
	staleKey := "stale-key"

	// Create two buckets
	manager.Allow(activeKey)
	manager.Allow(staleKey)

	// Should have 2 buckets
	if count := manager.Count(); count != 2 {
		t.Errorf("Expected 2 buckets, got %d", count)
	}

	// Keep using activeKey but not staleKey
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		for i := 0; i < 5; i++ {
			<-ticker.C
			manager.Allow(activeKey)
		}
		close(done)
	}()

	<-done

	// Wait for cleanup to run
	time.Sleep(cleanupInterval + 100*time.Millisecond)

	// Active key should remain, stale key should be cleaned
	if count := manager.Count(); count != 1 {
		t.Errorf("Expected 1 bucket remaining, got %d", count)
	}

	// Verify it's the active key that remains
	available := manager.AvailableTokens(activeKey)
	if available < 0 {
		t.Error("Active key bucket should still exist")
	}
}

func TestRateLimiterManager_Concurrent(t *testing.T) {
	manager := NewRateLimiterManager(100.0, 50.0, 1*time.Minute, 30*time.Minute)
	defer manager.Shutdown()

	numGoroutines := 20
	requestsPerGoroutine := 50

	var wg sync.WaitGroup

	// Multiple goroutines accessing different keys
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("user-%d", id)

			for j := 0; j < requestsPerGoroutine; j++ {
				manager.Allow(key)
			}
		}(i)
	}

	wg.Wait()

	// Should have created buckets for all goroutines
	count := manager.Count()
	if count != numGoroutines {
		t.Errorf("Expected %d buckets, got %d", numGoroutines, count)
	}
}

func TestRateLimiterManager_ConcurrentSameKey(t *testing.T) {
	manager := NewRateLimiterManager(100.0, 50.0, 1*time.Minute, 30*time.Minute)
	defer manager.Shutdown()

	numGoroutines := 50
	requestsPerGoroutine := 10
	key := "shared-key"

	var wg sync.WaitGroup
	allowedCount := make(chan int, numGoroutines)

	// Multiple goroutines accessing same key
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed := 0
			for j := 0; j < requestsPerGoroutine; j++ {
				if manager.Allow(key) {
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

	// Should allow approximately capacity requests
	if totalAllowed < 100 {
		t.Errorf("Expected at least 100 allowed requests, got %d", totalAllowed)
	}

	// Should not allow all requests (500 total)
	totalRequests := numGoroutines * requestsPerGoroutine
	if totalAllowed >= totalRequests {
		t.Errorf("Expected some requests to be denied, but all %d were allowed", totalRequests)
	}

	// Should only have 1 bucket for the shared key
	count := manager.Count()
	if count != 1 {
		t.Errorf("Expected 1 bucket for shared key, got %d", count)
	}

	t.Logf("Concurrent same-key test: %d/%d requests allowed", totalAllowed, totalRequests)
}

func TestRateLimiterManager_Shutdown(t *testing.T) {
	manager := NewRateLimiterManager(10.0, 5.0, 100*time.Millisecond, 1*time.Minute)

	// Create some buckets
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("key-%d", i)
		manager.Allow(key)
	}

	// Shutdown should complete without blocking
	done := make(chan bool)
	go func() {
		manager.Shutdown()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("Shutdown did not complete within timeout")
	}
}

func TestRateLimiterManager_MultipleKeys_IsolatedLimits(t *testing.T) {
	manager := NewRateLimiterManager(5.0, 1.0, 1*time.Minute, 30*time.Minute)
	defer manager.Shutdown()

	// Exhaust one key
	key1 := "ip:192.168.1.1"
	for i := 0; i < 5; i++ {
		manager.Allow(key1)
	}

	// Should be rate limited
	if manager.Allow(key1) {
		t.Error("Expected key1 to be rate limited")
	}

	// Different key should still work
	key2 := "ip:192.168.1.2"
	if !manager.Allow(key2) {
		t.Error("Expected key2 to be allowed (independent bucket)")
	}

	// Third key should also work
	key3 := "user:abc123"
	if !manager.Allow(key3) {
		t.Error("Expected key3 to be allowed (independent bucket)")
	}

	// Should have 3 independent buckets
	if count := manager.Count(); count != 3 {
		t.Errorf("Expected 3 buckets, got %d", count)
	}
}
