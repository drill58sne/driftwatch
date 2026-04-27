package filter_test

import (
	"testing"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/filter"
)

func sampleResults() []checker.CheckResult {
	return []checker.CheckResult{
		{Name: "cpu-usage", Host: "web-01", Output: "42%", Drift: false},
		{Name: "disk-usage", Host: "web-01", Output: "91%", Drift: true},
		{Name: "cpu-usage", Host: "db-01", Output: "88%", Drift: true},
		{Name: "memory", Host: "db-01", Output: "60%", Drift: false},
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	results := sampleResults()
	out := filter.Apply(results, filter.Options{})
	if len(out) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_OnlyDrift(t *testing.T) {
	out := filter.Apply(sampleResults(), filter.Options{OnlyDrift: true})
	if len(out) != 2 {
		t.Fatalf("expected 2 drifted results, got %d", len(out))
	}
	for _, r := range out {
		if !r.Drift {
			t.Errorf("expected Drift=true, got false for %s@%s", r.Name, r.Host)
		}
	}
}

func TestApply_FilterByTag(t *testing.T) {
	out := filter.Apply(sampleResults(), filter.Options{Tags: []string{"cpu"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 results matching 'cpu', got %d", len(out))
	}
	for _, r := range out {
		if r.Name != "cpu-usage" {
			t.Errorf("unexpected result name: %s", r.Name)
		}
	}
}

func TestApply_FilterByHost(t *testing.T) {
	out := filter.Apply(sampleResults(), filter.Options{Hosts: []string{"db-01"}})
	if len(out) != 2 {
		t.Fatalf("expected 2 results for db-01, got %d", len(out))
	}
	for _, r := range out {
		if r.Host != "db-01" {
			t.Errorf("unexpected host: %s", r.Host)
		}
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	out := filter.Apply(sampleResults(), filter.Options{
		OnlyDrift: true,
		Hosts:     []string{"web"},
	})
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Name != "disk-usage" || out[0].Host != "web-01" {
		t.Errorf("unexpected result: %+v", out[0])
	}
}

func TestApply_NoMatch_ReturnsNil(t *testing.T) {
	out := filter.Apply(sampleResults(), filter.Options{Hosts: []string{"nonexistent"}})
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d", len(out))
	}
}
