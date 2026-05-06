// Package throttle provides per-host connection throttling for driftwatch.
//
// When scanning many servers concurrently, opening too many SSH connections
// to a single host can trigger rate-limiting or exhaust server resources.
// Throttle enforces a configurable ceiling on simultaneous connections per
// host, blocking callers until a slot becomes available or the context is
// cancelled.
//
// Basic usage:
//
//	th, err := throttle.New(throttle.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if err := th.Acquire(ctx, host); err != nil {
//		return err
//	}
//	defer th.Release(host)
//
//	// open SSH connection and run checks …
package throttle
