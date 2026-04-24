package reporter

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func sampleResults() []DriftResult {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	return []DriftResult{
		{
			Host:      "web-01",
			CheckName: "nginx_version",
			Expected:  "1.24.0",
			Actual:    "1.24.0",
			Drifted:   false,
			Timestamp: now,
		},
		{
			Host:      "web-02",
			CheckName: "nginx_version",
			Expected:  "1.24.0",
			Actual:    "1.22.1",
			Drifted:   true,
			Timestamp: now,
		},
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	r := New(FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
	if r.format != FormatText {
		t.Errorf("expected format %q, got %q", FormatText, r.format)
	}
}

func TestWriteText_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	r := NewWithWriter(FormatText, &buf)
	if err := r.Write(sampleResults()); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	for _, col := range []string{"HOST", "CHECK", "STATUS", "EXPECTED", "ACTUAL"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected output to contain %q", col)
		}
	}
}

func TestWriteText_DriftStatus(t *testing.T) {
	var buf bytes.Buffer
	r := NewWithWriter(FormatText, &buf)
	r.Write(sampleResults())
	out := buf.String()
	if !strings.Contains(out, "DRIFT") {
		t.Error("expected DRIFT status in output")
	}
	if !strings.Contains(out, "OK") {
		t.Error("expected OK status in output")
	}
}

func TestWriteJSON_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	r := NewWithWriter(FormatJSON, &buf)
	if err := r.Write(sampleResults()); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(strings.TrimSpace(out), "[") {
		t.Error("expected JSON output to start with '['")
	}
	if !strings.Contains(out, `"drifted":true`) {
		t.Error("expected drifted:true in JSON output")
	}
	if !strings.Contains(out, `"host":"web-01"`) {
		t.Error("expected host web-01 in JSON output")
	}
}

func TestWriteText_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	r := NewWithWriter(FormatText, &buf)
	if err := r.Write([]DriftResult{}); err != nil {
		t.Fatalf("unexpected error on empty results: %v", err)
	}
}
