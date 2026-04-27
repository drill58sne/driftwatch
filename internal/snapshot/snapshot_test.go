package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/snapshot"
)

func sampleResults() []checker.CheckResult {
	return []checker.CheckResult{
		{Name: "os_version", Expected: "Ubuntu 22.04", Actual: "Ubuntu 22.04", Drift: false},
		{Name: "kernel", Expected: "5.15.0", Actual: "5.19.0", Drift: true},
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(dir)

	path, err := store.Save("web-01", sampleResults(), nil)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist at %s", path)
	}
}

func TestSave_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(dir)

	meta := map[string]string{"env": "prod"}
	path, err := store.Save("db-01", sampleResults(), meta)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	entry, err := store.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if entry.Host != "db-01" {
		t.Errorf("expected host db-01, got %s", entry.Host)
	}
	if len(entry.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(entry.Results))
	}
	if entry.Meta["env"] != "prod" {
		t.Errorf("expected meta env=prod, got %s", entry.Meta["env"])
	}
	if entry.CapturedAt.IsZero() {
		t.Error("expected non-zero CapturedAt")
	}
}

func TestSave_SanitizesHostname(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(dir)

	_, err := store.Save("192.168.1.1:22", sampleResults(), nil)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	entries, _ := filepath.Glob(filepath.Join(dir, "*.json"))
	if len(entries) != 1 {
		t.Fatalf("expected 1 snapshot file, got %d", len(entries))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	store := snapshot.NewStore(t.TempDir())
	_, err := store.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestEntry_CapturedAt_IsUTC(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewStore(dir)

	before := time.Now().UTC()
	path, _ := store.Save("host-x", sampleResults(), nil)
	after := time.Now().UTC()

	entry, err := store.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if entry.CapturedAt.Before(before) || entry.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v not in expected range [%v, %v]", entry.CapturedAt, before, after)
	}
}
