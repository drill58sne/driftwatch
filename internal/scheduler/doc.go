// Package scheduler provides a simple periodic job runner for driftwatch.
//
// It executes a drift-check Job on a configurable interval, running once
// immediately on start and then on each tick until the context is cancelled.
//
// Example usage:
//
//	opts := scheduler.Options{
//		Interval: 10 * time.Minute,
//		OnError: func(err error) { log.Println(err) },
//	}
//	s := scheduler.New(myJob, opts)
//	s.Run(ctx)
package scheduler
