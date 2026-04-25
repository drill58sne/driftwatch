// Package baseline provides functionality to save and load check result
// baselines, enabling drift comparison against a known-good state.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/driftwatch/internal/checker"
)

// Snapshot represents a saved baseline of check results.
type Snapshot struct {
	CreatedAt time.Time              `json:"created_at"`
	Host      string                 `json:"host"`
	Results   []checker.CheckResult  `json:"results"`
}

// Save writes a snapshot of results for the given host to the specified file path.
func Save(path, host string, results []checker.CheckResult) error {
	snap := Snapshot{
		CreatedAt: time.Now().UTC(),
		Host:      host,
		Results:   results,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal snapshot: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write file %q: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("baseline: read file %q: %w", path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal snapshot: %w", err)
	}
	return &snap, nil
}

// ResultMap converts a snapshot's results into a map keyed by check name
// for fast lookup during drift comparison.
func (s *Snapshot) ResultMap() map[string]checker.CheckResult {
	m := make(map[string]checker.CheckResult, len(s.Results))
	for _, r := range s.Results {
		m[r.Name] = r
	}
	return m
}
