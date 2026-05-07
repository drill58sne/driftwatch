package labeler_test

import (
	"testing"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/labeler"
)

func sampleResults() []checker.Result {
	return []checker.Result{
		{Host: "web-prod-01", Name: "nginx_version", Output: "1.24.0", Tags: []string{"service=nginx"}},
		{Host: "db-prod-01", Name: "postgres_version", Output: "15.2", Tags: nil},
		{Host: "web-staging-02", Name: "nginx_version", Output: "1.23.0", Tags: nil},
	}
}

func TestNew_ValidRules(t *testing.T) {
	_, err := labeler.New([]labeler.Rule{
		{Label: "env=prod", HostPattern: "prod"},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNew_InvalidHostPattern(t *testing.T) {
	_, err := labeler.New([]labeler.Rule{
		{Label: "env=prod", HostPattern: "["},
	})
	if err == nil {
		t.Fatal("expected error for invalid host pattern")
	}
}

func TestNew_InvalidCheckPattern(t *testing.T) {
	_, err := labeler.New([]labeler.Rule{
		{Label: "check=bad", CheckPattern: "("},
	})
	if err == nil {
		t.Fatal("expected error for invalid check pattern")
	}
}

func TestNew_EmptyLabel_ReturnsError(t *testing.T) {
	_, err := labeler.New([]labeler.Rule{
		{Label: "", HostPattern: "prod"},
	})
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestApply_MatchesHostPattern(t *testing.T) {
	l, _ := labeler.New([]labeler.Rule{
		{Label: "env=prod", HostPattern: "prod"},
	})
	out := l.Apply(sampleResults())
	for _, r := range out {
		if contains(r.Host, "prod") {
			if !hasTag(r.Tags, "env=prod") {
				t.Errorf("host %q missing label env=prod, tags=%v", r.Host, r.Tags)
			}
		} else {
			if hasTag(r.Tags, "env=prod") {
				t.Errorf("host %q should not have label env=prod", r.Host)
			}
		}
	}
}

func TestApply_MatchesCheckPattern(t *testing.T) {
	l, _ := labeler.New([]labeler.Rule{
		{Label: "check=nginx", CheckPattern: "nginx"},
	})
	out := l.Apply(sampleResults())
	for _, r := range out {
		if r.Name == "nginx_version" {
			if !hasTag(r.Tags, "check=nginx") {
				t.Errorf("result %q missing label check=nginx", r.Name)
			}
		}
	}
}

func TestApply_NoDuplicateLabels(t *testing.T) {
	l, _ := labeler.New([]labeler.Rule{
		{Label: "env=prod", HostPattern: "prod"},
		{Label: "env=prod", HostPattern: "web"},
	})
	out := l.Apply(sampleResults())
	for _, r := range out {
		count := 0
		for _, tag := range r.Tags {
			if tag == "env=prod" {
				count++
			}
		}
		if count > 1 {
			t.Errorf("host %q has duplicate label env=prod (%d times)", r.Host, count)
		}
	}
}

func TestApply_EmptyRules_ReturnsUnchanged(t *testing.T) {
	l, _ := labeler.New(nil)
	results := sampleResults()
	out := l.Apply(results)
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
}

func hasTag(tags []string, label string) bool {
	for _, t := range tags {
		if t == label {
			return true
		}
	}
	return false
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
