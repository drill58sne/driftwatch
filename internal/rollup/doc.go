// Package rollup provides aggregation of checker results into grouped
// summaries for reporting and alerting purposes.
//
// Results can be grouped by host or by tag. Each group entry exposes
// total, drifted, and clean counts, along with the raw results for
// downstream processing.
//
// Basic usage:
//
//	s := rollup.Aggregate(results, rollup.GroupByHost)
//	w := rollup.NewWriter(os.Stdout)
//	w.WriteText(s)
package rollup
