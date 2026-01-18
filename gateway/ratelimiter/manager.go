package ratelimiter

import (
	"sync"
	"time"
)

// RateLimiterManager manages multiple token buckets for different keys
type RateLimiterManager struct {
	limiters        sync.Map // map[string]*bucketWrapper
	capacity        float64
	refillRate      float64
	cleanupInterval time.Duration
	ttl             time.Duration
	stopCleanup     chan struct{}
	cleanupDone     chan struct{}
}

// bucketWrapper wraps a TokenBucket with metadata
type bucketWrapper struct {
	bucket   *TokenBucket
	lastUsed time.Time
	mu       sync.Mutex
}

// NewRateLimiterManager creates a new rate limiter manager
// capacity: maximim tokens per bucket (burst size)
// refillRate: tokens added per second per bucket
// cleanupInterval: how often to run cleanup
// ttl: time before removing the unused buckets
func NewRateLimiterManager(capacity, refillRate float64, cleanupInterval, ttl time.Duration) *RateLimiterManager {
	manager := &RateLimiterManager{
		capacity:        capacity,
		refillRate:      refillRate,
		cleanupInterval: cleanupInterval,
		ttl:             ttl,
		stopCleanup:     make(chan struct{}),
		cleanupDone:     make(chan struct{}),
	}

	// start background cleanup process
	go manager.cleanupLoop()

	return manager
}

// GetLimiter retrieves or creates a token bucket for the given key
func (m *RateLimiterManager) GetLimiter(key string) *TokenBucket {
	if value, ok := m.limiters.Load(key); ok {
		wrapper := value.(*bucketWrapper)
		wrapper.mu.Lock()
		wrapper.lastUsed = time.Now()
		wrapper.mu.Unlock()
		return wrapper.bucket
	}

	bucket := NewTokenBucket(m.capacity, m.refillRate)
	wrapper := &bucketWrapper{
		bucket:   bucket,
		lastUsed: time.Now(),
	}

	// Store it (LoadOrStore handles race condition)
	actual, loaded := m.limiters.LoadOrStore(key, wrapper)
	if loaded {
		wrapper = actual.(*bucketWrapper)
		wrapper.mu.Lock()
		wrapper.lastUsed = time.Now()
		wrapper.mu.Unlock()
		return wrapper.bucket
	}

	return bucket
}

// Allow checks if a request for the given key should be allowed
func (m *RateLimiterManager) Allow(key string) bool {
	bucket := m.GetLimiter(key)
	return bucket.Allow()
}

// AllowN checks if N requests for the given key should be allowed
func (m *RateLimiterManager) AllowN(key string, n int) bool {
	bucket := m.GetLimiter(key)
	return bucket.AllowN(n)
}

// AvailableTokens return the number of available tokens for a key
func (m *RateLimiterManager) AvailableTokens(key string) float64 {
	bucket := m.GetLimiter(key)
	return bucket.AvailableTokens()
}

// Shutdown stops the cleanup goroutine and releases resources
func (m *RateLimiterManager) Shutdown() {
	close(m.stopCleanup)
	<-m.cleanupDone
}

// cleanupLoop runs periodically to remove stale buckets
func (m *RateLimiterManager) cleanupLoop() {
	ticker := time.NewTicker(m.cleanupInterval)
	defer ticker.Stop()
	defer close(m.cleanupDone)

	for {
		select {
		case <-ticker.C:
			m.cleanup()
		case <-m.stopCleanup:
			return
		}
	}
}

// cleanup removes buckets that haven't been used within the TTL
func (m *RateLimiterManager) cleanup() {
	now := time.Now()
	keysToDelete := []string{}

	// collect keys to delete
	m.limiters.Range(func(key, value interface{}) bool {
		wrapper := value.(*bucketWrapper)
		wrapper.mu.Lock()
		lastUsed := wrapper.lastUsed
		wrapper.mu.Unlock()

		if now.Sub(lastUsed) > m.ttl {
			keysToDelete = append(keysToDelete, key.(string))
		}
		return true
	})

	// Delete stale buckets
	for _, key := range keysToDelete {
		m.limiters.Delete(key)
	}
}

// Count returns the number of active buckets (useful for testing/monitoring)
func (m *RateLimiterManager) Count() int {
	count := 0
	m.limiters.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}
