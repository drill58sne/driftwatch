// Package history provides functionality for recording and retrieving
// drift check run history across multiple executions of driftwatch.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/driftwatch/internal/differ"
)

// Entry represents a single recorded drift check run.
type Entry struct {
	Timestamp time.Time        `json:"timestamp"`
	Summary   differ.Summary   `json:"summary"`
	RunID     string           `json:"run_id"`
}

// Record holds all historical entries for a given history file.
type Record struct {
	Entries []Entry `json:"entries"`
}

// Append loads the history file at path, appends the new entry, and saves it.
// If the file does not exist it is created.
func Append(path string, summary differ.Summary) error {
	rec, err := load(path)
	if err != nil {
		return fmt.Errorf("history: load: %w", err)
	}

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Summary:   summary,
		RunID:     fmt.Sprintf("%d", time.Now().UnixNano()),
	}
	rec.Entries = append(rec.Entries, entry)

	if err := save(path, rec); err != nil {
		return fmt.Errorf("history: save: %w", err)
	}
	return nil
}

// Latest returns the most recent n entries from the history file at path,
// ordered newest-first. If n <= 0 all entries are returned.
func Latest(path string, n int) ([]Entry, error) {
	rec, err := load(path)
	if err != nil {
		return nil, fmt.Errorf("history: load: %w", err)
	}

	sort.Slice(rec.Entries, func(i, j int) bool {
		return rec.Entries[i].Timestamp.After(rec.Entries[j].Timestamp)
	})

	if n <= 0 || n >= len(rec.Entries) {
		return rec.Entries, nil
	}
	return rec.Entries[:n], nil
}

func load(path string) (Record, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Record{}, nil
	}
	if err != nil {
		return Record{}, err
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, err
	}
	return rec, nil
}

func save(path string, rec Record) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
