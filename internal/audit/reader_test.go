package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/driftwatch/internal/audit"
)

func TestReadAll_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	l, err := audit.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, host := range []string{"web-01", "db-01", "web-01"} {
		if err := l.Record(makeEvent(audit.EventScan, host, host == "db-01")); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}

	events, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 3 {
		t.Errorf("expected 3 events, got %d", len(events))
	}
}

func TestReadAll_FileNotFound_ReturnsNil(t *testing.T) {
	events, err := audit.ReadAll("/nonexistent/audit.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if events != nil {
		t.Errorf("expected nil events, got %v", events)
	}
}

func TestReadAll_CorruptLine_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.log")
	if err := os.WriteFile(path, []byte("not-json\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	_, err := audit.ReadAll(path)
	if err == nil {
		t.Error("expected error for corrupt line, got nil")
	}
}

func TestReadAll_EmptyFile_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.log")
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	events, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}
