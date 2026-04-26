// Package differ provides utilities for comparing checker results and
// identifying configuration drift between expected and actual system state.
//
// The primary entry point is Compare, which accepts a slice of checker.Result
// values and returns a differ.Result containing categorised drifted and clean
// entries. Use Summary to produce a human-readable overview suitable for
// CLI output or log messages.
//
// # Overview
//
// A typical usage pattern looks like:
//
//	results := checker.Run(ctx, checks)
//	diff := differ.Compare(results)
//	if diff.HasDrift() {
//		fmt.Println(differ.Summary(diff))
//	}
//
// Drifted entries are those where the observed state does not match the
// expected state defined in configuration. Clean entries confirm that the
// system is behaving as intended.
package differ
