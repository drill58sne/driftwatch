package timeout_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/timeout"
)

func TestDefaultOptions(t *testing.T) {
	opts := timeout.DefaultOptions()
	if opts.Dial <= 0 {
		t.Errorf("expected positive Dial, got %v", opts.Dial)
	}
	if opts.Exec <= 0 {
		t.Errorf("expected positive Exec, got %v", opts.Exec)
	}
	if opts.PerHost <= 0 {
		t.Errorf("expected positive PerHost, got %v", opts.PerHost)
	}
}

func TestValidate_Valid(t *testing.T) {
	opts := timeout.DefaultOptions()
	if err := opts.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_ZeroDial(t *testing.T) {
	opts := timeout.DefaultOptions()
	opts.Dial = 0
	if err := opts.Validate(); err == nil {
		t.Fatal("expected error for zero Dial, got nil")
	}
}

func TestValidate_ZeroExec(t *testing.T) {
	opts := timeout.DefaultOptions()
	opts.Exec = 0
	if err := opts.Validate(); err == nil {
		t.Fatal("expected error for zero Exec, got nil")
	}
}

func TestValidate_ZeroPerHost(t *testing.T) {
	opts := timeout.DefaultOptions()
	opts.PerHost = 0
	if err := opts.Validate(); err == nil {
		t.Fatal("expected error for zero PerHost, got nil")
	}
}

func TestForHost_CancelsAfterDeadline(t *testing.T) {
	opts := timeout.Options{
		Dial:    1 * time.Second,
		Exec:    1 * time.Second,
		PerHost: 50 * time.Millisecond,
	}
	ctx, cancel := opts.ForHost(context.Background())
	defer cancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context should have been cancelled within PerHost duration")
	}
}

func TestForExec_DeadlineSet(t *testing.T) {
	opts := timeout.DefaultOptions()
	ctx, cancel := opts.ForExec(context.Background())
	defer cancel()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	if time.Until(deadline) <= 0 {
		t.Fatal("expected deadline in the future")
	}
}

func TestForDial_DeadlineSet(t *testing.T) {
	opts := timeout.DefaultOptions()
	ctx, cancel := opts.ForDial(context.Background())
	defer cancel()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	if time.Until(deadline) <= 0 {
		t.Fatal("expected deadline in the future")
	}
}
