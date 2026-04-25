package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/baseline"
	"github.com/driftwatch/internal/checker"
)

func sampleResults() []checker.CheckResult {
	return []checker.CheckResult{
		{Name: "sshd_config", Expected: "PermitRootLogin no", Actual: "PermitRootLogin no", Drift: false},
		{Name: "timezone", Expected: "UTC", Actual: "America/New_York", Drift: true},
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	if err := baseline.Save(path, "host1", sampleResults()); err != nil {
		t.Fatalf("Save returned unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist at %q: %v", path, err)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	results := sampleResults()
	if err := baseline.Save(path, "web-01", results); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if snap.Host != "web-01" {
		t.Errorf("host: got %q, want %q", snap.Host, "web-01")
	}
	if len(snap.Results) != len(results) {
		t.Errorf("result count: got %d, want %d", len(snap.Results), len(results))
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := baseline.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestResultMap_KeyedByName(t *testing.T) {
	snap := &baseline.Snapshot{
		CreatedAt: time.Now(),
		Host:      "db-01",
		Results:   sampleResults(),
	}
	m := snap.ResultMap()
	if _, ok := m["sshd_config"]; !ok {
		t.Error("expected key 'sshd_config' in result map")
	}
	if _, ok := m["timezone"]; !ok {
		t.Error("expected key 'timezone' in result map")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o644)

	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
