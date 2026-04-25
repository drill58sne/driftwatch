// Package baseline manages the saving and loading of check result snapshots
// (baselines) for drift detection across remote hosts.
//
// A baseline is a point-in-time snapshot of check results for a given host.
// Once saved, subsequent runs can be compared against the baseline to detect
// configuration drift: values that have changed, checks that have been added,
// or checks that are no longer present.
//
// Usage:
//
//	// Save current results as the baseline
//	err := baseline.Save("/var/lib/driftwatch/web-01.json", "web-01", results)
//
//	// Load a previously saved baseline
//	snap, err := baseline.Load("/var/lib/driftwatch/web-01.json")
//
//	// Compare current results against the baseline
//	cr := snap.Against(currentResults)
//	fmt.Println(cr.Summary())
package baseline
