package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/differ"
	"github.com/user/driftwatch/internal/output"
)

func makeSummary(total, drifted, clean int) differ.Summary {
	return differ.Summary{
		Total:   total,
		Drifted: drifted,
		Clean:   clean,
	}
}

func TestNew_DefaultsToText(t *testing.T) {
	f := output.New(output.FormatText)
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestNewWithWriter_NilWriter_UsesStdout(t *testing.T) {
	f := output.NewWithWriter(output.FormatText, nil)
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestWrite_TextFormat_ContainsStatus(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewWithWriter(output.FormatText, &buf)

	if err := f.Write(makeSummary(5, 2, 3)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DRIFT DETECTED") {
		t.Errorf("expected 'DRIFT DETECTED' in output, got: %s", out)
	}
	if !strings.Contains(out, "Total") {
		t.Errorf("expected 'Total' label in output, got: %s", out)
	}
}

func TestWrite_TextFormat_CleanStatus(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewWithWriter(output.FormatText, &buf)

	if err := f.Write(makeSummary(3, 0, 3)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "clean") {
		t.Errorf("expected 'clean' status, got: %s", buf.String())
	}
}

func TestWrite_CompactFormat(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewWithWriter(output.FormatCompact, &buf)

	if err := f.Write(makeSummary(4, 1, 3)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, part := range []string{"total=4", "drifted=1", "clean=3"} {
		if !strings.Contains(out, part) {
			t.Errorf("expected %q in compact output, got: %s", part, out)
		}
	}
}

func TestWrite_JSONFormat_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewWithWriter(output.FormatJSON, &buf)

	if err := f.Write(makeSummary(6, 3, 3)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result differ.Summary
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if result.Total != 6 || result.Drifted != 3 || result.Clean != 3 {
		t.Errorf("unexpected values in JSON output: %+v", result)
	}
}
