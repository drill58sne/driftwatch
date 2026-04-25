package differ_test

import (
	"testing"

	"github.com/yourusername/driftwatch/internal/checker"
	"github.com/yourusername/driftwatch/internal/differ"
)

func sampleResults() []checker.Result {
	return []checker.Result{
		{Host: "web-01", Check: "os_version", Expected: "Ubuntu 22.04", Actual: "Ubuntu 20.04", Drift: true},
		{Host: "web-02", Check: "os_version", Expected: "Ubuntu 22.04", Actual: "Ubuntu 22.04", Drift: false},
		{Host: "db-01", Check: "kernel", Expected: "5.15", Actual: "5.10", Drift: true},
	}
}

func TestCompare_DetectsDrift(t *testing.T) {
	result := differ.Compare(sampleResults())
	if len(result.Drifted) != 2 {
		t.Errorf("expected 2 drifted, got %d", len(result.Drifted))
	}
}

func TestCompare_DetectsClean(t *testing.T) {
	result := differ.Compare(sampleResults())
	if len(result.Clean) != 1 {
		t.Errorf("expected 1 clean, got %d", len(result.Clean))
	}
}

func TestHasDrift_True(t *testing.T) {
	result := differ.Compare(sampleResults())
	if !result.HasDrift() {
		t.Error("expected HasDrift to be true")
	}
}

func TestHasDrift_False(t *testing.T) {
	clean := []checker.Result{
		{Host: "web-01", Check: "os_version", Expected: "Ubuntu 22.04", Actual: "Ubuntu 22.04", Drift: false},
	}
	result := differ.Compare(clean)
	if result.HasDrift() {
		t.Error("expected HasDrift to be false")
	}
}

func TestSummary_NoDrift(t *testing.T) {
	result := &differ.Result{}
	s := differ.Summary(result)
	if s != "No drift detected." {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithDrift(t *testing.T) {
	result := differ.Compare(sampleResults())
	s := differ.Summary(result)
	if s == "" {
		t.Error("expected non-empty summary for drifted result")
	}
}

func TestDiff_Fields(t *testing.T) {
	result := differ.Compare(sampleResults())
	d := result.Drifted[0]
	if d.Host == "" || d.Check == "" || d.Expected == "" || d.Actual == "" {
		t.Error("expected all Diff fields to be populated")
	}
}
