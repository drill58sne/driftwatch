package redactor_test

import (
	"testing"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/redactor"
)

func sampleResults() []checker.CheckResult {
	return []checker.CheckResult{
		{Name: "os_version", Output: "Ubuntu 22.04", Expected: "Ubuntu 22.04", Drift: false},
		{Name: "api_key", Output: "s3cr3t-key-value", Expected: "s3cr3t-key-value", Drift: false},
		{Name: "db_password", Output: "hunter2", Expected: "hunter2", Drift: false},
		{Name: "token_count", Output: "42", Expected: "42", Drift: false},
	}
}

func TestNew_DefaultPatterns(t *testing.T) {
	r, err := redactor.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil redactor")
	}
}

func TestNew_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := redactor.New([]string{`[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestApply_MasksSensitiveFields(t *testing.T) {
	r, _ := redactor.New(nil)
	results := sampleResults()
	got := r.Apply(results)

	sensitive := map[string]bool{"api_key": true, "db_password": true, "token_count": true}
	for _, res := range got {
		if sensitive[res.Name] {
			if res.Output != "[REDACTED]" {
				t.Errorf("name=%q: expected Output=[REDACTED], got %q", res.Name, res.Output)
			}
			if res.Expected != "[REDACTED]" {
				t.Errorf("name=%q: expected Expected=[REDACTED], got %q", res.Name, res.Expected)
			}
		}
	}
}

func TestApply_PreservesNonSensitiveFields(t *testing.T) {
	r, _ := redactor.New(nil)
	results := sampleResults()
	got := r.Apply(results)

	for _, res := range got {
		if res.Name == "os_version" {
			if res.Output != "Ubuntu 22.04" {
				t.Errorf("expected Output to be unchanged, got %q", res.Output)
			}
			return
		}
	}
	t.Error("os_version result not found")
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	r, _ := redactor.New(nil)
	results := sampleResults()
	r.Apply(results)

	for _, res := range results {
		if res.Name == "api_key" && res.Output == "[REDACTED]" {
			t.Error("original slice was mutated")
		}
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	r, err := redactor.New([]string{`(?i)os_version`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	results := sampleResults()
	got := r.Apply(results)

	for _, res := range got {
		if res.Name == "os_version" && res.Output != "[REDACTED]" {
			t.Errorf("expected os_version to be redacted with custom pattern")
		}
		if res.Name == "api_key" && res.Output == "[REDACTED]" {
			t.Errorf("api_key should not be redacted with custom pattern")
		}
	}
}
