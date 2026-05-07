package rollup_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/rollup"
)

func buildSummary() rollup.Summary {
	return rollup.Aggregate(sampleResults(), rollup.GroupByHost)
}

func TestNewWriter_NilUsesStdout(t *testing.T) {
	w := rollup.NewWriter(nil)
	if w == nil {
		t.Fatal("expected non-nil Writer")
	}
}

func TestWriteText_ContainsGroupByLabel(t *testing.T) {
	var buf bytes.Buffer
	w := rollup.NewWriter(&buf)
	s := buildSummary()
	if err := w.WriteText(s); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	if !strings.Contains(buf.String(), "host") {
		t.Error("expected 'host' in text output")
	}
}

func TestWriteText_ContainsDriftStatus(t *testing.T) {
	var buf bytes.Buffer
	w := rollup.NewWriter(&buf)
	if err := w.WriteText(buildSummary()); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	if !strings.Contains(buf.String(), "DRIFT") {
		t.Error("expected DRIFT label in output")
	}
}

func TestWriteJSON_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	w := rollup.NewWriter(&buf)
	if err := w.WriteJSON(buildSummary()); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["Entries"]; !ok {
		t.Error("expected 'Entries' key in JSON output")
	}
}

func TestWriteText_AllHostsPresent(t *testing.T) {
	var buf bytes.Buffer
	w := rollup.NewWriter(&buf)
	w.WriteText(buildSummary())
	out := buf.String()
	for _, host := range []string{"web-01", "db-01", "cache-01"} {
		if !strings.Contains(out, host) {
			t.Errorf("expected host %q in output", host)
		}
	}
}
