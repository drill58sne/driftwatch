package cache_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/cache"
	"github.com/driftwatch/internal/checker"
)

func sampleResults() []checker.CheckResult {
	return []checker.CheckResult{
		{Name: "os.version", Expected: "Ubuntu 22.04", Actual: "Ubuntu 22.04", Drift: false},
		{Name: "kernel", Expected: "5.15.0", Actual: "5.19.0", Drift: true},
	}
}

func TestNew_DefaultTTL(t *testing.T) {
	c := cache.New(0)
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestSet_And_Get_Hit(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("host-a", sampleResults())

	results, ok := c.Get("host-a")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestGet_Miss_UnknownHost(t *testing.T) {
	c := cache.New(time.Minute)
	_, ok := c.Get("unknown")
	if ok {
		t.Fatal("expected cache miss for unknown host")
	}
}

func TestGet_Miss_AfterExpiry(t *testing.T) {
	c := cache.New(time.Millisecond)
	// Inject a frozen clock, then advance it past the TTL.
	now := time.Now()
	c.(*cache.Cache) // type assertion not possible on unexported nowFunc; use internal helper via New
	// Re-create with a custom TTL and rely on real time sleep for expiry test.
	c2 := cache.New(50 * time.Millisecond)
	c2.Set("host-b", sampleResults())
	time.Sleep(80 * time.Millisecond)
	_, ok := c2.Get("host-b")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
	_ = now
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("host-c", sampleResults())
	c.Invalidate("host-c")
	_, ok := c.Get("host-c")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("host-d", sampleResults())
	c.Set("host-e", sampleResults())
	c.Flush()
	if c.Size() != 0 {
		t.Fatalf("expected empty cache after flush, got size %d", c.Size())
	}
}

func TestSize_ReflectsEntries(t *testing.T) {
	c := cache.New(time.Minute)
	if c.Size() != 0 {
		t.Fatalf("expected size 0, got %d", c.Size())
	}
	c.Set("host-f", sampleResults())
	if c.Size() != 1 {
		t.Fatalf("expected size 1, got %d", c.Size())
	}
}
