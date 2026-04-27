// Package filter provides utilities for filtering check results
// based on tags, host patterns, and drift status.
package filter

import (
	"strings"

	"github.com/driftwatch/internal/checker"
)

// Options holds filtering criteria for check results.
type Options struct {
	// OnlyDrift filters results to only those with drift detected.
	OnlyDrift bool
	// Tags filters results to only those whose check name contains one of the given tags.
	Tags []string
	// Hosts filters results to only those matching one of the given host substrings.
	Hosts []string
}

// Apply filters a slice of CheckResult according to the provided Options.
// It returns a new slice containing only the results that match all criteria.
func Apply(results []checker.CheckResult, opts Options) []checker.CheckResult {
	var out []checker.CheckResult
	for _, r := range results {
		if opts.OnlyDrift && !r.Drift {
			continue
		}
		if len(opts.Tags) > 0 && !matchesAny(r.Name, opts.Tags) {
			continue
		}
		if len(opts.Hosts) > 0 && !matchesAny(r.Host, opts.Hosts) {
			continue
		}
		out = append(out, r)
	}
	return out
}

// matchesAny returns true if value contains any of the given substrings.
func matchesAny(value string, patterns []string) bool {
	for _, p := range patterns {
		if strings.Contains(value, p) {
			return true
		}
	}
	return false
}
