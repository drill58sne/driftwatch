// Package ratelimit provides a token-bucket rate limiter for controlling
// the frequency of SSH connections and check executions across hosts.
package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Options configures the rate limiter.
type Options struct {
	// Rate is the number of tokens added per second.
	Rate float64
	// Burst is the maximum number of tokens the bucket can hold.
	Burst int
}

// DefaultOptions returns sensible defaults for the rate limiter.
func DefaultOptions() Options {
	return Options{
		Rate:  5.0,
		Burst: 10,
	}
}

// Limiter controls the rate of operations using a token bucket.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	burst    float64
	rate     float64
	lastTick time.Time
}

// New creates a new Limiter with the given options.
// Returns an error if options are invalid.
func New(opts Options) (*Limiter, error) {
	if opts.Rate <= 0 {
		return nil, fmt.Errorf("ratelimit: rate must be positive, got %f", opts.Rate)
	}
	if opts.Burst <= 0 {
		return nil, fmt.Errorf("ratelimit: burst must be positive, got %d", opts.Burst)
	}
	return &Limiter{
		tokens:   float64(opts.Burst),
		burst:    float64(opts.Burst),
		rate:     opts.Rate,
		lastTick: time.Now(),
	}, nil
}

// Wait blocks until a token is available or the context is cancelled.
func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("ratelimit: context cancelled: %w", err)
		}
		if l.tryConsume() {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("ratelimit: context cancelled: %w", ctx.Err())
		case <-time.After(time.Duration(float64(time.Second) / l.rate)):
		}
	}
}

// tryConsume attempts to consume one token. Returns true if successful.
func (l *Limiter) tryConsume() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.burst {
		l.tokens = l.burst
	}

	if l.tokens >= 1.0 {
		l.tokens -= 1.0
		return true
	}
	return false
}
