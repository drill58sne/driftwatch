package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/differ"
	"github.com/driftwatch/internal/history"
)

func makeSummary(total, drifted int) differ.Summary {
	return differ.Summary{
		Total:   total,
		Drifted: drifted,
		Clean:   total - drifted,
	}
}

func TestAppend_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	if err := history.Append(path, makeSummary(3, 1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestAppend_AccumulatesEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	for i := 0; i < 3; i++ {
		if err := history.Append(path, makeSummary(5, i)); err != nil {
			t.Fatalf("append %d failed: %v", i, err)
		}
	}

	entries, err := history.Latest(path, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestLatest_ReturnsNewestFirst(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	for i := 0; i < 4; i++ {
		time.Sleep(time.Millisecond)
		if err := history.Append(path, makeSummary(4, i)); err != nil {
			t.Fatalf("append failed: %v", err)
		}
	}

	entries, err := history.Latest(path, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if !entries[0].Timestamp.After(entries[1].Timestamp) {
		t.Errorf("expected entries ordered newest-first")
	}
}

func TestLatest_FileNotFound_ReturnsEmpty(t *testing.T) {
	entries, err := history.Latest("/nonexistent/path/history.json", 0)
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty entries for missing file, got %d", len(entries))
	}
}

func TestEntry_HasRunID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	if err := history.Append(path, makeSummary(2, 0)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := history.Latest(path, 1)
	if entries[0].RunID == "" {
		t.Error("expected non-empty RunID")
	}
}
