package tagger_test

import (
	"testing"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/tagger"
)

// TestGroupThenKeys exercises the full Group→Keys pipeline to ensure the two
// functions stay consistent with each other across a realistic dataset.
func TestGroupThenKeys_Consistent(t *testing.T) {
	results := []checker.CheckResult{
		{Host: "a", Name: "check", Tags: []string{"env=prod", "region=us-east"}},
		{Host: "b", Name: "check", Tags: []string{"env=staging", "region=eu-west"}},
		{Host: "c", Name: "check", Tags: []string{"env=prod", "region=us-east"}},
	}

	keys := tagger.Keys(results)
	for _, key := range keys {
		groups := tagger.Group(results, key)
		total := 0
		for _, g := range groups {
			total += len(g)
		}
		if total != len(results) {
			t.Errorf("key %q: grouped %d results, want %d", key, total, len(results))
		}
	}
}

// TestGroup_LargeDataset ensures grouping scales without panicking.
func TestGroup_LargeDataset(t *testing.T) {
	envs := []string{"prod", "staging", "dev", "qa"}
	results := make([]checker.CheckResult, 0, 400)
	for i := 0; i < 400; i++ {
		env := envs[i%len(envs)]
		results = append(results, checker.CheckResult{
			Host:  "host",
			Name:  "check",
			Tags:  []string{"env=" + env},
		})
	}

	groups := tagger.Group(results, "env")
	if len(groups) != len(envs) {
		t.Errorf("got %d groups, want %d", len(groups), len(envs))
	}
	for _, env := range envs {
		if len(groups[env]) != 100 {
			t.Errorf("env=%s: got %d, want 100", env, len(groups[env]))
		}
	}
}

// TestParseTag_RoundTrip verifies that String() output can be re-parsed.
func TestParseTag_RoundTrip(t *testing.T) {
	original := tagger.Tag{Key: "team", Value: "platform"}
	parsed, err := tagger.ParseTag(original.String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed != original {
		t.Errorf("round-trip mismatch: got %+v, want %+v", parsed, original)
	}
}
