// Package differ provides utilities for comparing checker results and
// identifying configuration drift between expected and actual system state.
//
// The primary entry point is Compare, which accepts a slice of checker.Result
// values and returns a differ.Result containing categorised drifted and clean
// entries. Use Summary to produce a human-readable overview suitable for
// CLI output or log messages.
package differ
