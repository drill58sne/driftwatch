package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/alert"
	"github.com/user/driftwatch/internal/differ"
)

func makeSummary(total, drifted int) differ.Summary {
	return differ.Summary{
		Total:   total,
		Drifted: drifted,
		Clean:   total - drifted,
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := alert.DefaultConfig()
	if cfg.WarnThreshold <= 0 {
		t.Error("expected positive WarnThreshold")
	}
	if cfg.ErrorThreshold <= cfg.WarnThreshold {
		t.Error("expected ErrorThreshold > WarnThreshold")
	}
}

func TestEvaluate_NoDrift(t *testing.T) {
	a := alert.New(alert.DefaultConfig())
	result := a.Evaluate(makeSummary(5, 0))
	if result != nil {
		t.Errorf("expected nil alert for no drift, got %+v", result)
	}
}

func TestEvaluate_WarnLevel(t *testing.T) {
	cfg := alert.Config{WarnThreshold: 1, ErrorThreshold: 5}
	a := alert.New(cfg)
	result := a.Evaluate(makeSummary(5, 2))
	if result == nil {
		t.Fatal("expected an alert")
	}
	if result.Level != alert.LevelWarn {
		t.Errorf("expected WARN, got %s", result.Level)
	}
	if result.Count != 2 {
		t.Errorf("expected count 2, got %d", result.Count)
	}
}

func TestEvaluate_ErrorLevel(t *testing.T) {
	cfg := alert.Config{WarnThreshold: 1, ErrorThreshold: 5}
	a := alert.New(cfg)
	result := a.Evaluate(makeSummary(10, 6))
	if result == nil {
		t.Fatal("expected an alert")
	}
	if result.Level != alert.LevelError {
		t.Errorf("expected ERROR, got %s", result.Level)
	}
}

func TestEmit_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	cfg := alert.Config{WarnThreshold: 1, ErrorThreshold: 5}
	a := alert.NewWithWriter(cfg, &buf)
	result := a.Evaluate(makeSummary(3, 2))
	a.Emit(result)

	out := buf.String()
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected [WARN] in output, got: %s", out)
	}
	if !strings.Contains(out, "2 drifted") {
		t.Errorf("expected drift count in output, got: %s", out)
	}
}

func TestEmit_NilAlert_NoOutput(t *testing.T) {
	var buf bytes.Buffer
	a := alert.NewWithWriter(alert.DefaultConfig(), &buf)
	a.Emit(nil)
	if buf.Len() != 0 {
		t.Errorf("expected no output for nil alert, got: %s", buf.String())
	}
}
