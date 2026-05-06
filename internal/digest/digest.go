// Package digest provides utilities for computing and comparing
// content hashes of check results to detect meaningful changes
// between successive drift scans.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/driftwatch/internal/checker"
)

// Result holds the computed digest for a single check result.
type Result struct {
	Host   string
	Name   string
	Digest string
}

// Compute returns a SHA-256 hex digest of the given CheckResult's
// output and status, ignoring fields that do not affect drift.
func Compute(r checker.CheckResult) (Result, error) {
	payload := struct {
		Status string
		Output string
	}{
		Status: r.Status,
		Output: r.Output,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return Result{}, fmt.Errorf("digest: marshal: %w", err)
	}

	sum := sha256.Sum256(b)
	return Result{
		Host:   r.Host,
		Name:   r.Name,
		Digest: hex.EncodeToString(sum[:]),
	}, nil
}

// ComputeAll returns a digest Result for every entry in the slice.
// The first error encountered halts processing.
func ComputeAll(results []checker.CheckResult) ([]Result, error) {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		d, err := Compute(r)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

// Changed returns true when the digest of current differs from previous
// for the same host+name pair.  If the pair is not found in previous,
// Changed returns true (treated as a new / changed result).
func Changed(previous []Result, current Result) bool {
	for _, p := range previous {
		if p.Host == current.Host && p.Name == current.Name {
			return p.Digest != current.Digest
		}
	}
	return true
}
