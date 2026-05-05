package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/retry"
)

func TestDefaultOptions(t *testing.T) {
	opts := retry.DefaultOptions()
	if opts.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", opts.MaxAttempts)
	}
	if opts.InitialDelay != 500*time.Millisecond {
		t.Errorf("unexpected InitialDelay: %v", opts.InitialDelay)
	}
	if opts.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", opts.Multiplier)
	}
}

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.DefaultOptions(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	calls := 0
	sentinel := errors.New("transient")
	opts := retry.Options{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond, Multiplier: 2}

	err := retry.Do(context.Background(), opts, func() error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success on third attempt, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	calls := 0
	opts := retry.Options{MaxAttempts: 2, InitialDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond, Multiplier: 1}

	err := retry.Do(context.Background(), opts, func() error {
		calls++
		return errors.New("always fails")
	})
	if err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestDo_RespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := retry.Do(ctx, retry.DefaultOptions(), func() error {
		return nil
	})
	if err == nil {
		t.Fatal("expected error due to cancelled context")
	}
}

func TestDo_InvalidMaxAttempts(t *testing.T) {
	opts := retry.Options{MaxAttempts: 0}
	err := retry.Do(context.Background(), opts, func() error { return nil })
	if err == nil {
		t.Fatal("expected error for MaxAttempts=0")
	}
}
