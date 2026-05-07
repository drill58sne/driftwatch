package rollup_test

import (
	"testing"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/rollup"
)

func sampleResults() []checker.Result {
	return []checker.Result{
		{Host: "web-01", Name: "nginx", Drift: true, Tags: []string{"web"}},
		{Host: "web-01", Name: "sshd", Drift: false, Tags: []string{"web"}},
		{Host: "db-01", Name: "postgres", Drift: false, Tags: []string{"db"}},
		{Host: "db-01", Name: "sshd", Drift: true, Tags: []string{"db"}},
		{Host: "cache-01", Name: "redis", Drift: false, Tags: []string{}},
	}
}

func TestAggregate_GroupByHost_KeyCount(t *testing.T) {
	s := rollup.Aggregate(sampleResults(), rollup.GroupByHost)
	if len(s.Entries) != 3 {
		t.Fatalf("expected 3 host entries, got %d", len(s.Entries))
	}
}

func TestAggregate_GroupByHost_Counts(t *testing.T) {
	s := rollup.Aggregate(sampleResults(), rollup.GroupByHost)
	for _, e := range s.Entries {
		if e.Total != e.Drifted+e.Clean {
			t.Errorf("host %s: total %d != drifted %d + clean %d", e.Key, e.Total, e.Drifted, e.Clean)
		}
	}
}

func TestAggregate_GroupByTag_UntaggedFallback(t *testing.T) {
	s := rollup.Aggregate(sampleResults(), rollup.GroupByTag)
	var found bool
	for _, e := range s.Entries {
		if e.Key == "untagged" {
			found = true
		}
	}
	if !found {
		t.Error("expected an 'untagged' entry for result with no tags")
	}
}

func TestEntry_HasDrift_True(t *testing.T) {
	e := rollup.Entry{Drifted: 2, Clean: 1}
	if !e.HasDrift() {
		t.Error("expected HasDrift true")
	}
}

func TestEntry_HasDrift_False(t *testing.T) {
	e := rollup.Entry{Drifted: 0, Clean: 3}
	if e.HasDrift() {
		t.Error("expected HasDrift false")
	}
}

func TestAggregate_EmptyResults(t *testing.T) {
	s := rollup.Aggregate(nil, rollup.GroupByHost)
	if len(s.Entries) != 0 {
		t.Errorf("expected 0 entries for empty input, got %d", len(s.Entries))
	}
}

func TestAggregate_SortedByKey(t *testing.T) {
	s := rollup.Aggregate(sampleResults(), rollup.GroupByHost)
	for i := 1; i < len(s.Entries); i++ {
		if s.Entries[i].Key < s.Entries[i-1].Key {
			t.Errorf("entries not sorted: %s before %s", s.Entries[i-1].Key, s.Entries[i].Key)
		}
	}
}
