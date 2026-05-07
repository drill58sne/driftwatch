package summarizer_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftwatch/internal/summarizer"
)

func buildSummary() summarizer.Summary {
	return summarizer.Compute(makeResults())
}

func TestNewWriter_NilUsesStdout(t *testing.T) {
	w := summarizer.NewWriter(nil, "text")
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestWriteText_ContainsSummaryLine(t *testing.T) {
	var buf bytes.Buffer
	w := summarizer.NewWriter(&buf, "text")
	if err := w.Write(buildSummary()); err != nil {
		t.Fatalf("Write() error: %v", err)
	}
	if !strings.Contains(buf.String(), "total=") {
		t.Errorf("expected 'total=' in output, got: %s", buf.String())
	}
}

func TestWriteText_DriftedStatus(t *testing.T) {
	var buf bytes.Buffer
	w := summarizer.NewWriter(&buf, "text")
	_ = w.Write(buildSummary())
	if !strings.Contains(buf.String(), "DRIFTED") {
		t.Errorf("expected 'DRIFTED' in output, got: %s", buf.String())
	}
}

func TestWriteJSON_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	w := summarizer.NewWriter(&buf, "json")
	if err := w.Write(buildSummary()); err != nil {
		t.Fatalf("Write() error: %v", err)
	}
	var out summarizer.Summary
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Total != 4 {
		t.Errorf("expected Total=4, got %d", out.Total)
	}
}

func TestWriteText_UnknownFormatDefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	w := summarizer.NewWriter(&buf, "xml")
	_ = w.Write(buildSummary())
	if !strings.Contains(buf.String(), "Summary:") {
		t.Errorf("expected text output for unknown format, got: %s", buf.String())
	}
}
