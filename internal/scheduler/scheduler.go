// Package scheduler provides periodic drift-check execution support.
package scheduler

import (
	"context"
	"log"
	"time"
)

// Job represents a function to be executed on a schedule.
type Job func(ctx context.Context) error

// Options configures the scheduler.
type Options struct {
	// Interval between job executions.
	Interval time.Duration
	// OnError is called when a job returns an error. If nil, errors are logged.
	OnError func(err error)
}

// DefaultOptions returns sensible scheduler defaults.
func DefaultOptions() Options {
	return Options{
		Interval: 5 * time.Minute,
		OnError: func(err error) {
			log.Printf("[scheduler] job error: %v", err)
		},
	}
}

// Scheduler runs a Job repeatedly on a fixed interval.
type Scheduler struct {
	opts Options
	job  Job
}

// New creates a new Scheduler with the given job and options.
func New(job Job, opts Options) *Scheduler {
	if opts.OnError == nil {
		opts.OnError = DefaultOptions().OnError
	}
	return &Scheduler{opts: opts, job: job}
}

// Run starts the scheduler loop, blocking until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	log.Printf("[scheduler] starting, interval=%s", s.opts.Interval)
	ticker := time.NewTicker(s.opts.Interval)
	defer ticker.Stop()

	// Run once immediately before waiting for first tick.
	if err := s.job(ctx); err != nil {
		s.opts.OnError(err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("[scheduler] stopped")
			return
		case <-ticker.C:
			if err := s.job(ctx); err != nil {
				s.opts.OnError(err)
			}
		}
	}
}
