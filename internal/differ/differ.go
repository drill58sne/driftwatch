// Package differ compares expected vs actual check results to identify drift.
package differ

import "fmt"

import "github.com/yourusername/driftwatch/internal/checker"

// Diff represents a single drift difference between expected and actual values.
type Diff struct {
	Host     string
	Check    string
	Expected string
	Actual   string
}

// Result holds the full diff output for a set of check results.
type Result struct {
	Drifted []Diff
	Clean   []string // host:check pairs with no drift
}

// HasDrift returns true if any drifted entries exist.
func (r *Result) HasDrift() bool {
	return len(r.Drifted) > 0
}

// Compare takes a slice of checker results and separates drifted from clean.
func Compare(results []checker.Result) *Result {
	out := &Result{}

	for _, r := range results {
		key := r.Host + ":" + r.Check
		if r.Drift {
			out.Drifted = append(out.Drifted, Diff{
				Host:     r.Host,
				Check:    r.Check,
				Expected: r.Expected,
				Actual:   r.Actual,
			})
		} else {
			out.Clean = append(out.Clean, key)
		}
	}

	return out
}

// Summary returns a human-readable summary string.
func Summary(r *Result) string {
	if !r.HasDrift() {
		return "No drift detected."
	}
	return fmt.Sprintf("%d drift(s) detected across %d clean check(s).",
		len(r.Drifted), len(r.Clean))
}
