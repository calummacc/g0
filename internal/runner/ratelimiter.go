package runner

import (
	"context"
	"time"
)

// RateLimiter implements a token bucket rate limiter
// It ensures that requests don't exceed the specified rate per second
type RateLimiter struct {
	tokens   chan struct{}
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewRateLimiter creates a new rate limiter with the specified max RPS
// If maxRPS is 0 or negative, rate limiting is disabled (returns nil)
func NewRateLimiter(maxRPS int) *RateLimiter {
	if maxRPS <= 0 {
		return nil // No rate limiting
	}

	ctx, cancel := context.WithCancel(context.Background())
	rl := &RateLimiter{
		tokens:   make(chan struct{}, maxRPS), // Buffer allows burst up to maxRPS
		interval: time.Second / time.Duration(maxRPS),
		ctx:      ctx,
		cancel:   cancel,
	}

	// Pre-fill the bucket with tokens
	for i := 0; i < maxRPS; i++ {
		rl.tokens <- struct{}{}
	}

	// Start token refill goroutine
	go rl.refill()

	return rl
}

// refill continuously adds tokens to the bucket at the specified rate
func (rl *RateLimiter) refill() {
	ticker := time.NewTicker(rl.interval)
	defer ticker.Stop()

	for {
		select {
		case <-rl.ctx.Done():
			return
		case <-ticker.C:
			// Try to add a token, but don't block if bucket is full
			select {
			case rl.tokens <- struct{}{}:
			default:
				// Bucket is full, skip
			}
		}
	}
}

// Wait blocks until a token is available, ensuring rate limit is respected
// Returns false if context is cancelled
func (rl *RateLimiter) Wait(ctx context.Context) bool {
	if rl == nil {
		return true // No rate limiting, proceed immediately
	}

	select {
	case <-ctx.Done():
		return false
	case <-rl.ctx.Done():
		return false
	case <-rl.tokens:
		return true // Token acquired, proceed
	}
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
	if rl != nil {
		rl.cancel()
	}
}

