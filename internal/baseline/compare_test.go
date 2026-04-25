package baseline_test

import (
	"strings"
	"testing"
	"time"

	"github.com/driftwatch/internal/baseline"
	"github.com/driftwatch/internal/checker"
)

func makeSnapshot(results []checker.CheckResult) *baseline.Snapshot {
	return &baseline.Snapshot{
		CreatedAt: time.Now(),
		Host:      "web-01",
		Results:   results,
	}
}

func TestAgainst_NoDrift(t *testing.T) {
	results := sampleResults()
	snap := makeSnapshot(results)
	cr := snap.Against(results)
	if cr.HasDrift() {
		t.Errorf("expected no drift, got: %+v", cr)
	}
}

func TestAgainst_DetectsDrifted(t *testing.T) {
	base := []checker.CheckResult{
		{Name: "timezone", Expected: "UTC", Actual: "UTC", Drift: false},
	}
	current := []checker.CheckResult{
		{Name: "timezone", Expected: "UTC", Actual: "America/New_York", Drift: true},
	}
	snap := makeSnapshot(base)
	cr := snap.Against(current)
	if len(cr.Drifted) != 1 {
		t.Fatalf("expected 1 drifted entry, got %d", len(cr.Drifted))
	}
	if cr.Drifted[0].Name != "timezone" {
		t.Errorf("drifted name: got %q, want %q", cr.Drifted[0].Name, "timezone")
	}
}

func TestAgainst_DetectsNew(t *testing.T) {
	snap := makeSnapshot([]checker.CheckResult{})
	current := []checker.CheckResult{
		{Name: "new_check", Expected: "on", Actual: "on", Drift: false},
	}
	cr := snap.Against(current)
	if len(cr.New) != 1 {
		t.Fatalf("expected 1 new entry, got %d", len(cr.New))
	}
}

func TestAgainst_DetectsRemoved(t *testing.T) {
	base := []checker.CheckResult{
		{Name: "old_check", Expected: "yes", Actual: "yes", Drift: false},
	}
	snap := makeSnapshot(base)
	cr := snap.Against([]checker.CheckResult{})
	if len(cr.Removed) != 1 {
		t.Fatalf("expected 1 removed entry, got %d", len(cr.Removed))
	}
	if cr.Removed[0] != "old_check" {
		t.Errorf("removed name: got %q, want %q", cr.Removed[0], "old_check")
	}
}

func TestSummary_NoDrift(t *testing.T) {
	snap := makeSnapshot(sampleResults())
	cr := snap.Against(sampleResults())
	if !strings.Contains(cr.Summary(), "no drift") {
		t.Errorf("summary should mention 'no drift', got: %q", cr.Summary())
	}
}

func TestSummary_WithDrift(t *testing.T) {
	snap := makeSnapshot([]checker.CheckResult{
		{Name: "tz", Actual: "UTC"},
	})
	cr := snap.Against([]checker.CheckResult{
		{Name: "tz", Actual: "EST"},
	})
	if !strings.Contains(cr.Summary(), "drifted") {
		t.Errorf("summary should mention 'drifted', got: %q", cr.Summary())
	}
}
