package ratelimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   float64    // max number of token in the bucket
	tokens     float64    // current available tokens
	refillRate float64    // tokens added per second
	lastRefill time.Time  // last time that tokens were refilled
	mu         sync.Mutex // concurrent access
}

func NewTokenBucket(capacity, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity, // start with full bucket
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request can be allowed and consumes 1 token
// Returns true if allowed, false if rate limited
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}

	return false
}

// AllowN checks if N tokens are available and consumes them
// Returns true if allowed, else false if rate limited
func (tb *TokenBucket) AllowN(n int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	requiredTokens := float64(n)
	if tb.tokens >= requiredTokens {
		tb.tokens -= requiredTokens
		return true
	}
	return false
}

// refill adds tokens based on elapsed time since last refill
// Must be called with lock held
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()

	// calculate tokens to add based on elapsed time
	tokensToAdd := elapsed * tb.refillRate

	// add tokens but don't exceed capacity
	tb.tokens = min(tb.tokens+tokensToAdd, tb.capacity)

	tb.lastRefill = now
}

func (tb *TokenBucket) AllowWithRemaining() (allowed bool, remaining float64) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0

		return true, tb.tokens
	}

	return false, tb.tokens
}

// AvailableTokens returns the current number of available tokens
// Useful for debugging and rate limit headers
func (tb *TokenBucket) AvailableTokens() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	return tb.tokens
}
