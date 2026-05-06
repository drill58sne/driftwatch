// Package throttle provides a host-level SSH connection throttle that limits
// the number of concurrent connections opened to any single remote host.
package throttle

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrThrottled is returned when a host slot cannot be acquired.
var ErrThrottled = errors.New("throttle: max concurrent connections reached for host")

// Options configures the Throttle.
type Options struct {
	// MaxPerHost is the maximum number of simultaneous connections allowed per host.
	MaxPerHost int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxPerHost: 3,
	}
}

// Throttle tracks per-host connection counts and enforces a ceiling.
type Throttle struct {
	mu         sync.Mutex
	maxPerHost int
	counts     map[string]int
	waiters    map[string][]chan struct{}
}

// New creates a Throttle with the given options.
// Returns an error if MaxPerHost is less than 1.
func New(opts Options) (*Throttle, error) {
	if opts.MaxPerHost < 1 {
		return nil, fmt.Errorf("throttle: MaxPerHost must be >= 1, got %d", opts.MaxPerHost)
	}
	return &Throttle{
		maxPerHost: opts.MaxPerHost,
		counts:     make(map[string]int),
		waiters:    make(map[string][]chan struct{}),
	}, nil
}

// Acquire blocks until a connection slot is available for host or ctx is done.
// Callers must call Release when the connection is closed.
func (t *Throttle) Acquire(ctx context.Context, host string) error {
	for {
		t.mu.Lock()
		if t.counts[host] < t.maxPerHost {
			t.counts[host]++
			t.mu.Unlock()
			return nil
		}
		ch := make(chan struct{}, 1)
		t.waiters[host] = append(t.waiters[host], ch)
		t.mu.Unlock()

		select {
		case <-ctx.Done():
			t.removeWaiter(host, ch)
			return fmt.Errorf("throttle: %w: %s", ctx.Err(), host)
		case <-ch:
			// slot may now be free; re-check at top of loop
		}
	}
}

// Release frees a connection slot for host.
func (t *Throttle) Release(host string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.counts[host] > 0 {
		t.counts[host]--
	}
	if len(t.waiters[host]) > 0 {
		ch := t.waiters[host][0]
		t.waiters[host] = t.waiters[host][1:]
		ch <- struct{}{}
	}
}

// Available returns the number of free slots for host.
func (t *Throttle) Available(host string) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.maxPerHost - t.counts[host]
}

func (t *Throttle) removeWaiter(host string, ch chan struct{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	list := t.waiters[host]
	for i, c := range list {
		if c == ch {
			t.waiters[host] = append(list[:i], list[i+1:]...)
			return
		}
	}
}
