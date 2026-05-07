// Package rollup aggregates multiple check results into a single
// consolidated summary, grouping by host or tag for reporting.
package rollup

import (
	"sort"
	"time"

	"github.com/driftwatch/internal/checker"
)

// GroupBy defines how results should be aggregated.
type GroupBy string

const (
	GroupByHost GroupBy = "host"
	GroupByTag  GroupBy = "tag"
)

// Entry holds aggregated results for a single group key.
type Entry struct {
	Key       string
	Total     int
	Drifted   int
	Clean     int
	Results   []checker.Result
	UpdatedAt time.Time
}

// HasDrift returns true if any result in the entry has drifted.
func (e Entry) HasDrift() bool {
	return e.Drifted > 0
}

// Summary is the full rollup output, indexed by group key.
type Summary struct {
	GroupBy GroupBy
	Entries []Entry
	RolledAt time.Time
}

// Aggregate groups the provided results by the given strategy and
// returns a Summary with per-group counts.
func Aggregate(results []checker.Result, by GroupBy) Summary {
	groups := make(map[string][]checker.Result)

	for _, r := range results {
		key := keyFor(r, by)
		groups[key] = append(groups[key], r)
	}

	entries := make([]Entry, 0, len(groups))
	for key, rs := range groups {
		e := Entry{
			Key:       key,
			Total:     len(rs),
			Results:   rs,
			UpdatedAt: time.Now(),
		}
		for _, r := range rs {
			if r.Drift {
				e.Drifted++
			} else {
				e.Clean++
			}
		}
		entries = append(entries, e)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Summary{
		GroupBy:  by,
		Entries:  entries,
		RolledAt: time.Now(),
	}
}

func keyFor(r checker.Result, by GroupBy) string {
	switch by {
	case GroupByTag:
		if len(r.Tags) > 0 {
			return r.Tags[0]
		}
		return "untagged"
	default:
		return r.Host
	}
}
