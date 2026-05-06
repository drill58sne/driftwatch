package throttle_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/throttle"
)

func TestDefaultOptions(t *testing.T) {
	opts := throttle.DefaultOptions()
	if opts.MaxPerHost != 3 {
		t.Fatalf("expected MaxPerHost=3, got %d", opts.MaxPerHost)
	}
}

func TestNew_InvalidMaxPerHost(t *testing.T) {
	_, err := throttle.New(throttle.Options{MaxPerHost: 0})
	if err == nil {
		t.Fatal("expected error for MaxPerHost=0")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	th, err := throttle.New(throttle.Options{MaxPerHost: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil Throttle")
	}
}

func TestAcquire_Release_UpdatesAvailable(t *testing.T) {
	th, _ := throttle.New(throttle.Options{MaxPerHost: 2})
	ctx := context.Background()

	if got := th.Available("host-a"); got != 2 {
		t.Fatalf("expected 2 available, got %d", got)
	}
	if err := th.Acquire(ctx, "host-a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := th.Available("host-a"); got != 1 {
		t.Fatalf("expected 1 available after acquire, got %d", got)
	}
	th.Release("host-a")
	if got := th.Available("host-a"); got != 2 {
		t.Fatalf("expected 2 available after release, got %d", got)
	}
}

func TestAcquire_BlocksWhenFull(t *testing.T) {
	th, _ := throttle.New(throttle.Options{MaxPerHost: 1})
	ctx := context.Background()

	// Fill the single slot.
	if err := th.Acquire(ctx, "host-b"); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}

	acquired := make(chan struct{})
	go func() {
		_ = th.Acquire(ctx, "host-b")
		close(acquired)
	}()

	select {
	case <-acquired:
		t.Fatal("second acquire should have blocked")
	case <-time.After(50 * time.Millisecond):
	}

	th.Release("host-b")
	select {
	case <-acquired:
		// success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("second acquire did not unblock after release")
	}
}

func TestAcquire_ContextCancellation(t *testing.T) {
	th, _ := throttle.New(throttle.Options{MaxPerHost: 1})
	ctx, cancel := context.WithCancel(context.Background())

	_ = th.Acquire(context.Background(), "host-c")

	errCh := make(chan error, 1)
	go func() {
		errCh <- th.Acquire(ctx, "host-c")
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected error on cancelled context")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("acquire did not return after context cancel")
	}
}

func TestAcquire_ConcurrentHosts_Independent(t *testing.T) {
	th, _ := throttle.New(throttle.Options{MaxPerHost: 1})
	ctx := context.Background()
	var wg sync.WaitGroup
	hosts := []string{"h1", "h2", "h3"}
	for _, h := range hosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			if err := th.Acquire(ctx, host); err != nil {
				t.Errorf("acquire %s: %v", host, err)
			}
			time.Sleep(10 * time.Millisecond)
			th.Release(host)
		}(h)
	}
	wg.Wait()
}
