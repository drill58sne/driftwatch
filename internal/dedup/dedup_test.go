package dedup_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/dedup"
)

func makeResult(host, name, output string, drift bool) checker.CheckResult {
	return checker.CheckResult{
		Host:   host,
		Name:   name,
		Output: output,
		Drift:  drift,
	}
}

func TestFilter_NewResult_PassesThrough(t *testing.T) {
	s := dedup.New(time.Minute)
	results := []checker.CheckResult{makeResult("host1", "check1", "v1", true)}
	out := s.Filter(results)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
}

func TestFilter_DuplicateResult_Suppressed(t *testing.T) {
	s := dedup.New(time.Minute)
	r := makeResult("host1", "check1", "v1", true)
	s.Filter([]checker.CheckResult{r})
	out := s.Filter([]checker.CheckResult{r})
	if len(out) != 0 {
		t.Fatalf("expected 0 results, got %d", len(out))
	}
}

func TestFilter_ChangedOutput_PassesThrough(t *testing.T) {
	s := dedup.New(time.Minute)
	s.Filter([]checker.CheckResult{makeResult("host1", "check1", "v1", true)})
	out := s.Filter([]checker.CheckResult{makeResult("host1", "check1", "v2", true)})
	if len(out) != 1 {
		t.Fatalf("expected 1 result after output change, got %d", len(out))
	}
}

func TestFilter_AfterTTLExpiry_PassesThrough(t *testing.T) {
	s := dedup.New(1 * time.Millisecond)
	r := makeResult("host1", "check1", "v1", true)
	s.Filter([]checker.CheckResult{r})
	time.Sleep(5 * time.Millisecond)
	out := s.Filter([]checker.CheckResult{r})
	if len(out) != 1 {
		t.Fatalf("expected result to pass through after TTL, got %d", len(out))
	}
}

func TestStats_TracksSuppressedCount(t *testing.T) {
	s := dedup.New(time.Minute)
	r := makeResult("host1", "check1", "v1", true)
	s.Filter([]checker.CheckResult{r})
	s.Filter([]checker.CheckResult{r})
	s.Filter([]checker.CheckResult{r})
	_, suppressed := s.Stats()
	if suppressed != 2 {
		t.Fatalf("expected 2 suppressed, got %d", suppressed)
	}
}

func TestEvict_RemovesExpiredEntries(t *testing.T) {
	s := dedup.New(1 * time.Millisecond)
	s.Filter([]checker.CheckResult{makeResult("host1", "check1", "v1", true)})
	time.Sleep(5 * time.Millisecond)
	s.Evict()
	entries, _ := s.Stats()
	if entries != 0 {
		t.Fatalf("expected 0 entries after eviction, got %d", entries)
	}
}

func TestNew_DefaultTTL_WhenZero(t *testing.T) {
	s := dedup.New(0)
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}
