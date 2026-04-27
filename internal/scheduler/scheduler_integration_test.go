package scheduler_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/scheduler"
)

// TestRun_TicksMultipleTimes verifies the scheduler fires more than once
// when the interval is short enough to tick within the test window.
func TestRun_TicksMultipleTimes(t *testing.T) {
	var count int64
	job := func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		return nil
	}

	opts := scheduler.Options{
		Interval: 30 * time.Millisecond,
		OnError: func(err error) { t.Errorf("unexpected error: %v", err) },
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	s := scheduler.New(job, opts)
	s.Run(ctx)

	final := atomic.LoadInt64(&count)
	if final < 2 {
		t.Errorf("expected at least 2 executions, got %d", final)
	}
}

// TestRun_StopsOnContextCancel verifies the scheduler stops promptly.
func TestRun_StopsOnContextCancel(t *testing.T) {
	var count int64
	job := func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		return nil
	}

	opts := scheduler.Options{
		Interval: 5 * time.Millisecond,
		OnError: func(err error) { t.Errorf("unexpected error: %v", err) },
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		s := scheduler.New(job, opts)
		s.Run(ctx)
		close(done)
	}()

	time.Sleep(25 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// success
	case <-time.After(200 * time.Millisecond):
		t.Error("scheduler did not stop after context cancellation")
	}
}
