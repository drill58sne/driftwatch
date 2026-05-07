package tagger_test

import (
	"testing"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/tagger"
)

func sampleResults() []checker.CheckResult {
	return []checker.CheckResult{
		{Host: "web-01", Name: "uptime", Tags: []string{"env=prod", "team=platform"}},
		{Host: "web-02", Name: "uptime", Tags: []string{"env=prod", "team=sre"}},
		{Host: "db-01", Name: "uptime", Tags: []string{"env=staging", "team=platform"}},
		{Host: "bare", Name: "uptime", Tags: []string{}},
	}
}

func TestParseTag_Valid(t *testing.T) {
	tag, err := tagger.ParseTag("env=prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Key != "env" || tag.Value != "prod" {
		t.Errorf("got %+v, want {env prod}", tag)
	}
}

func TestParseTag_Invalid(t *testing.T) {
	cases := []string{"noequals", "=value", "key=", ""}
	for _, c := range cases {
		_, err := tagger.ParseTag(c)
		if err == nil {
			t.Errorf("expected error for %q, got nil", c)
		}
	}
}

func TestTag_String(t *testing.T) {
	tag := tagger.Tag{Key: "env", Value: "prod"}
	if tag.String() != "env=prod" {
		t.Errorf("got %q, want \"env=prod\"", tag.String())
	}
}

func TestGroup_ByEnv(t *testing.T) {
	results := sampleResults()
	groups := tagger.Group(results, "env")

	if len(groups["prod"]) != 2 {
		t.Errorf("prod group: got %d, want 2", len(groups["prod"]))
	}
	if len(groups["staging"]) != 1 {
		t.Errorf("staging group: got %d, want 1", len(groups["staging"]))
	}
	if len(groups["(untagged)"]) != 1 {
		t.Errorf("untagged group: got %d, want 1", len(groups["(untagged)"]))
	}
}

func TestGroup_ByTeam(t *testing.T) {
	results := sampleResults()
	groups := tagger.Group(results, "team")

	if len(groups["platform"]) != 2 {
		t.Errorf("platform group: got %d, want 2", len(groups["platform"]))
	}
	if len(groups["sre"]) != 1 {
		t.Errorf("sre group: got %d, want 1", len(groups["sre"]))
	}
}

func TestGroup_UnknownKey_AllUntagged(t *testing.T) {
	results := sampleResults()
	groups := tagger.Group(results, "region")
	if len(groups["(untagged)"]) != len(results) {
		t.Errorf("expected all %d results untagged, got %d", len(results), len(groups["(untagged)"]))
	}
}

func TestKeys_ReturnsSortedUniqueKeys(t *testing.T) {
	results := sampleResults()
	keys := tagger.Keys(results)

	if len(keys) != 2 {
		t.Fatalf("got %d keys, want 2: %v", len(keys), keys)
	}
	if keys[0] != "env" || keys[1] != "team" {
		t.Errorf("got %v, want [env team]", keys)
	}
}

func TestKeys_EmptyResults(t *testing.T) {
	keys := tagger.Keys(nil)
	if len(keys) != 0 {
		t.Errorf("expected empty keys, got %v", keys)
	}
}
