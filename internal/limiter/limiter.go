// Package limiter provides a concurrency limiter for bounding the number
// of simultaneous SSH connections or check executions across hosts.
package limiter

import (
	"context"
	"errors"
)

// ErrLimitExceeded is returned when Acquire is called with a cancelled context.
var ErrLimitExceeded = errors.New("limiter: context cancelled while waiting to acquire slot")

// Limiter bounds the number of concurrent operations using a semaphore channel.
type Limiter struct {
	sem chan struct{}
}

// Options configures the Limiter.
type Options struct {
	// MaxConcurrent is the maximum number of concurrent operations allowed.
	// Must be >= 1.
	MaxConcurrent int
}

// DefaultOptions returns sensible defaults for the Limiter.
func DefaultOptions() Options {
	return Options{
		MaxConcurrent: 10,
	}
}

// New creates a new Limiter with the given options.
// If MaxConcurrent is less than 1, it defaults to 1.
func New(opts Options) *Limiter {
	if opts.MaxConcurrent < 1 {
		opts.MaxConcurrent = 1
	}
	return &Limiter{
		sem: make(chan struct{}, opts.MaxConcurrent),
	}
}

// Acquire blocks until a slot is available or the context is cancelled.
// Callers must call Release after the operation completes.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case l.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ErrLimitExceeded
	}
}

// Release frees a previously acquired slot.
func (l *Limiter) Release() {
	<-l.sem
}

// Available returns the number of slots currently free.
func (l *Limiter) Available() int {
	return cap(l.sem) - len(l.sem)
}

// Capacity returns the total number of concurrent slots.
func (l *Limiter) Capacity() int {
	return cap(l.sem)
}
