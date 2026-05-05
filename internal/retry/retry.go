// Package retry provides configurable retry logic with exponential backoff
// for use when connecting to remote hosts or executing checks over SSH.
package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Options configures the retry behaviour.
type Options struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// InitialDelay is the wait time before the second attempt.
	InitialDelay time.Duration
	// MaxDelay caps the exponential backoff ceiling.
	MaxDelay time.Duration
	// Multiplier is the factor applied to the delay on each failure.
	Multiplier float64
}

// DefaultOptions returns sensible defaults suitable for SSH operations.
func DefaultOptions() Options {
	return Options{
		MaxAttempts:  3,
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// Do executes fn up to opts.MaxAttempts times, backing off between failures.
// The context is checked before every attempt; cancellation stops retries
// immediately and returns the context error.
func Do(ctx context.Context, opts Options, fn func() error) error {
	if opts.MaxAttempts <= 0 {
		return errors.New("retry: MaxAttempts must be greater than zero")
	}
	if opts.Multiplier <= 0 {
		opts.Multiplier = 1
	}

	delay := opts.InitialDelay
	var lastErr error

	for attempt := 1; attempt <= opts.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry: context cancelled before attempt %d: %w", attempt, err)
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if attempt == opts.MaxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("retry: context cancelled during backoff: %w", ctx.Err())
		case <-time.After(delay):
		}

		delay = time.Duration(float64(delay) * opts.Multiplier)
		if delay > opts.MaxDelay {
			delay = opts.MaxDelay
		}
	}

	return fmt.Errorf("retry: all %d attempts failed: %w", opts.MaxAttempts, lastErr)
}
