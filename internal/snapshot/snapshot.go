// Package snapshot provides functionality for capturing and storing
// point-in-time snapshots of check results across all hosts.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/driftwatch/internal/checker"
)

// Entry represents a single snapshot captured at a point in time.
type Entry struct {
	CapturedAt time.Time                       `json:"captured_at"`
	Host       string                          `json:"host"`
	Results    []checker.CheckResult           `json:"results"`
	Meta       map[string]string               `json:"meta,omitempty"`
}

// Store manages snapshot persistence on disk.
type Store struct {
	dir string
}

// NewStore returns a Store that persists snapshots under dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Save writes a snapshot entry for the given host to disk.
// Files are named <host>_<unix-nano>.json inside the store directory.
func (s *Store) Save(host string, results []checker.CheckResult, meta map[string]string) (string, error) {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return "", fmt.Errorf("snapshot: mkdir %s: %w", s.dir, err)
	}

	entry := Entry{
		CapturedAt: time.Now().UTC(),
		Host:       host,
		Results:    results,
		Meta:       meta,
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return "", fmt.Errorf("snapshot: marshal: %w", err)
	}

	filename := fmt.Sprintf("%s_%d.json", sanitize(host), entry.CapturedAt.UnixNano())
	path := filepath.Join(s.dir, filename)

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("snapshot: write %s: %w", path, err)
	}

	return path, nil
}

// Load reads a snapshot entry from the given file path.
func (s *Store) Load(path string) (Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Entry{}, fmt.Errorf("snapshot: read %s: %w", path, err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return Entry{}, fmt.Errorf("snapshot: unmarshal %s: %w", path, err)
	}

	return entry, nil
}

// sanitize replaces characters unsuitable for filenames with underscores.
func sanitize(s string) string {
	out := make([]byte, len(s))
	for i := range s {
		if s[i] == '/' || s[i] == ':' || s[i] == '\\' {
			out[i] = '_'
		} else {
			out[i] = s[i]
		}
	}
	return string(out)
}
