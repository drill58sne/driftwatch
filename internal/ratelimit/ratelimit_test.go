package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/driftwatch/internal/ratelimit"
)

func TestDefaultOptions(t *testing.T) {
	opts := ratelimit.DefaultOptions()
	if opts.Rate <= 0 {
		t.Errorf("expected positive rate, got %f", opts.Rate)
	}
	if opts.Burst <= 0 {
		t.Errorf("expected positive burst, got %d", opts.Burst)
	}
}

func TestNew_InvalidRate(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Options{Rate: -1, Burst: 5})
	if err == nil {
		t.Error("expected error for negative rate")
	}
}

func TestNew_InvalidBurst(t *testing.T) {
	_, err := ratelimit.New(ratelimit.Options{Rate: 1.0, Burst: 0})
	if err == nil {
		t.Error("expected error for zero burst")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Options{Rate: 10.0, Burst: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestWait_ConsumesTokens(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Options{Rate: 100.0, Burst: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	// Should succeed immediately since burst=5
	for i := 0; i < 5; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("Wait() failed on iteration %d: %v", i, err)
		}
	}
}

func TestWait_ContextCancelled(t *testing.T) {
	// Very low rate so tokens run out quickly
	l, err := ratelimit.New(ratelimit.Options{Rate: 0.01, Burst: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	// Consume the single burst token
	_ = l.Wait(ctx)

	// Next call should block; cancel the context
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = l.Wait(ctx)
	if err == nil {
		t.Error("expected error due to context cancellation")
	}
}

func TestWait_HighRate_IsNonBlocking(t *testing.T) {
	l, err := ratelimit.New(ratelimit.Options{Rate: 1000.0, Burst: 50})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	start := time.Now()
	for i := 0; i < 10; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("Wait() failed: %v", err)
		}
	}
	if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
		t.Errorf("expected fast execution, took %v", elapsed)
	}
}
