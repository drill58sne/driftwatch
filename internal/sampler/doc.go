// Package sampler provides a sliding-window sample store for checker results.
//
// It is designed for lightweight in-memory trend analysis: results are recorded
// per host and automatically evicted when they exceed the configured window size
// or maximum age.
//
// Basic usage:
//
//	s := sampler.New(sampler.DefaultOptions())
//	s.Record("web-01", results)
//	window := s.Get("web-01")
//
The zero value is not usable; always construct via New.
package sampler
