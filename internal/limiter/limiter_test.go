package limiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/driftwatch/internal/limiter"
)

func TestDefaultOptions(t *testing.T) {
	opts := limiter.DefaultOptions()
	if opts.MaxConcurrent != 10 {
		t.Errorf("expected MaxConcurrent=10, got %d", opts.MaxConcurrent)
	}
}

func TestNew_ClampsBelowOne(t *testing.T) {
	l := limiter.New(limiter.Options{MaxConcurrent: 0})
	if l.Capacity() != 1 {
		t.Errorf("expected capacity 1 for zero value, got %d", l.Capacity())
	}
}

func TestNew_SetsCapacity(t *testing.T) {
	l := limiter.New(limiter.Options{MaxConcurrent: 5})
	if l.Capacity() != 5 {
		t.Errorf("expected capacity 5, got %d", l.Capacity())
	}
}

func TestAcquire_Release_UpdatesAvailable(t *testing.T) {
	l := limiter.New(limiter.Options{MaxConcurrent: 3})

	if l.Available() != 3 {
		t.Fatalf("expected 3 available before acquire, got %d", l.Available())
	}

	_ = l.Acquire(context.Background())
	if l.Available() != 2 {
		t.Errorf("expected 2 available after one acquire, got %d", l.Available())
	}

	l.Release()
	if l.Available() != 3 {
		t.Errorf("expected 3 available after release, got %d", l.Available())
	}
}

func TestAcquire_BlocksWhenFull(t *testing.T) {
	l := limiter.New(limiter.Options{MaxConcurrent: 1})
	_ = l.Acquire(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := l.Acquire(ctx)
	if err == nil {
		t.Error("expected error when limiter is full and context expires")
	}
}

func TestAcquire_CancelledContext_ReturnsError(t *testing.T) {
	l := limiter.New(limiter.Options{MaxConcurrent: 1})
	_ = l.Acquire(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := l.Acquire(ctx)
	if err != limiter.ErrLimitExceeded {
		t.Errorf("expected ErrLimitExceeded, got %v", err)
	}
}

func TestAcquire_ConcurrentGoroutines_RespectsLimit(t *testing.T) {
	const max = 3
	l := limiter.New(limiter.Options{MaxConcurrent: max})

	var mu sync.Mutex
	peak := 0
	current := 0
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = l.Acquire(context.Background())
			defer l.Release()

			mu.Lock()
			current++
			if current > peak {
				peak = current
			}
			mu.Unlock()

			time.Sleep(10 * time.Millisecond)

			mu.Lock()
			current--
			mu.Unlock()
		}()
	}

	wg.Wait()

	if peak > max {
		t.Errorf("peak concurrency %d exceeded limit %d", peak, max)
	}
}
