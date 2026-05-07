// Package tagger provides utilities for parsing, grouping, and inspecting
// tags attached to check results.
//
// Tags follow the "key=value" convention and are stored as plain strings on
// checker.CheckResult.Tags.  This package does not mutate results; it only
// reads the tag slice and organises results into labelled buckets.
//
// Typical usage:
//
//	groups := tagger.Group(results, "env")
//	for env, rs := range groups {
//	    fmt.Printf("env=%s  results=%d\n", env, len(rs))
//	}
//
// Use Keys to discover which tag keys are present before grouping:
//
//	for _, k := range tagger.Keys(results) {
//	    fmt.Println("tag key:", k)
//	}
package tagger
