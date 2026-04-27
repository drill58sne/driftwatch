// Package snapshot captures and persists point-in-time snapshots of
// check results for one or more remote hosts.
//
// A snapshot records the full set of [checker.CheckResult] values returned
// during a single drift-check run, together with the host name, a UTC
// timestamp, and optional metadata key/value pairs (e.g. environment tags).
//
// Snapshots are stored as JSON files on disk via [Store]. Each file is named
// using the sanitized hostname and a nanosecond Unix timestamp to ensure
// uniqueness across rapid successive captures.
//
// Typical usage:
//
//	store := snapshot.NewStore("/var/lib/driftwatch/snapshots")
//	path, err := store.Save(host, results, map[string]string{"env": "prod"})
//	if err != nil { ... }
//
//	entry, err := store.Load(path)
//	if err != nil { ... }
package snapshot
