package scheduler_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/scheduler"
)

func TestDefaultOptions(t *testing.T) {
	opts := scheduler.DefaultOptions()
	if opts.Interval != 5*time.Minute {
		t.Errorf("expected 5m interval, got %s", opts.Interval)
	}
	if opts.OnError == nil {
		t.Error("expected non-nil OnError handler")
	}
}

func TestNew_NotNil(t *testing.T) {
	s := scheduler.New(func(ctx context.Context) error { return nil }, scheduler.DefaultOptions())
	if s == nil {
		t.Fatal("expected non-nil Scheduler")
	}
}

func TestRun_ExecutesJobImmediately(t *testing.T) {
	var count int64
	job := func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		return nil
	}

	opts := scheduler.Options{
		Interval: 10 * time.Second, // long interval so only immediate run fires
		OnError: func(err error) { t.Errorf("unexpected error: %v", err) },
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	s := scheduler.New(job, opts)
	s.Run(ctx)

	if atomic.LoadInt64(&count) < 1 {
		t.Error("expected job to run at least once immediately")
	}
}

func TestRun_CallsOnError(t *testing.T) {
	var errCount int64
	job := func(ctx context.Context) error {
		return errors.New("job failed")
	}

	opts := scheduler.Options{
		Interval: 10 * time.Second,
		OnError: func(err error) {
			atomic.AddInt64(&errCount, 1)
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	s := scheduler.New(job, opts)
	s.Run(ctx)

	if atomic.LoadInt64(&errCount) < 1 {
		t.Error("expected OnError to be called at least once")
	}
}

func TestRun_NilOnError_UsesDefault(t *testing.T) {
	// Should not panic when OnError is nil (scheduler provides default).
	job := func(ctx context.Context) error {
		return errors.New("some error")
	}

	opts := scheduler.Options{
		Interval: 10 * time.Second,
		OnError:  nil,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	s := scheduler.New(job, opts)
	s.Run(ctx) // should not panic
}
