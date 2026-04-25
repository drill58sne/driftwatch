package baseline

import (
	"fmt"

	"github.com/driftwatch/internal/checker"
)

// DriftEntry describes a single check that has changed relative to the baseline.
type DriftEntry struct {
	Name     string
	Baseline string
	Current  string
}

// CompareResult holds the outcome of comparing current results against a baseline.
type CompareResult struct {
	Host     string
	Drifted  []DriftEntry
	New      []checker.CheckResult // checks present in current but not in baseline
	Removed  []string              // check names present in baseline but not current
}

// HasDrift reports whether any drift, new, or removed checks were detected.
func (c *CompareResult) HasDrift() bool {
	return len(c.Drifted) > 0 || len(c.New) > 0 || len(c.Removed) > 0
}

// Summary returns a human-readable one-line summary of the comparison.
func (c *CompareResult) Summary() string {
	if !c.HasDrift() {
		return fmt.Sprintf("host %s: no drift detected", c.Host)
	}
	return fmt.Sprintf("host %s: %d drifted, %d new, %d removed",
		c.Host, len(c.Drifted), len(c.New), len(c.Removed))
}

// Against compares current check results against the snapshot baseline and
// returns a CompareResult describing what has changed.
func (s *Snapshot) Against(current []checker.CheckResult) *CompareResult {
	baseMap := s.ResultMap()
	currMap := make(map[string]checker.CheckResult, len(current))
	for _, r := range current {
		currMap[r.Name] = r
	}

	out := &CompareResult{Host: s.Host}

	for _, r := range current {
		b, exists := baseMap[r.Name]
		if !exists {
			out.New = append(out.New, r)
			continue
		}
		if r.Actual != b.Actual {
			out.Drifted = append(out.Drifted, DriftEntry{
				Name:     r.Name,
				Baseline: b.Actual,
				Current:  r.Actual,
			})
		}
	}

	for name := range baseMap {
		if _, exists := currMap[name]; !exists {
			out.Removed = append(out.Removed, name)
		}
	}

	return out
}
