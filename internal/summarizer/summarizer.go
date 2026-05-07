// Package summarizer aggregates check results into a concise per-host summary
// suitable for display, alerting, or further processing.
package summarizer

import (
	"time"

	"github.com/driftwatch/internal/checker"
)

// HostSummary holds aggregated drift statistics for a single host.
type HostSummary struct {
	Host       string
	Total      int
	Drifted    int
	Clean      int
	Errored    int
	DriftRate  float64 // fraction of checks that drifted
	CheckedAt  time.Time
}

// Summary is the full aggregation across all hosts.
type Summary struct {
	Hosts     []HostSummary
	Total     int
	Drifted   int
	Clean     int
	Errored   int
	CreatedAt time.Time
}

// HasDrift reports whether any host in the summary has drifted checks.
func (s Summary) HasDrift() bool {
	return s.Drifted > 0
}

// Compute builds a Summary from a slice of checker results.
func Compute(results []checker.Result) Summary {
	hostMap := make(map[string]*HostSummary)

	for _, r := range results {
		hs, ok := hostMap[r.Host]
		if !ok {
			hs = &HostSummary{Host: r.Host, CheckedAt: r.CheckedAt}
			hostMap[r.Host] = hs
		}
		hs.Total++
		switch {
		case r.Error != nil:
			hs.Errored++
		case r.Drifted:
			hs.Drifted++
		default:
			hs.Clean++
		}
	}

	summary := Summary{CreatedAt: time.Now()}
	for _, hs := range hostMap {
		if hs.Total > 0 {
			hs.DriftRate = float64(hs.Drifted) / float64(hs.Total)
		}
		summary.Hosts = append(summary.Hosts, *hs)
		summary.Total += hs.Total
		summary.Drifted += hs.Drifted
		summary.Clean += hs.Clean
		summary.Errored += hs.Errored
	}
	return summary
}
