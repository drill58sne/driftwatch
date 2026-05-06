// Package timeout provides per-host SSH operation timeout management
// for driftwatch check runs.
package timeout

import (
	"context"
	"fmt"
	"time"
)

// DefaultOptions returns a Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Dial:      10 * time.Second,
		Exec:      30 * time.Second,
		PerHost:   60 * time.Second,
	}
}

// Options holds timeout durations for each phase of a check.
type Options struct {
	// Dial is the maximum time allowed to establish an SSH connection.
	Dial time.Duration
	// Exec is the maximum time allowed for a single remote command.
	Exec time.Duration
	// PerHost is the total budget for all checks against one host.
	PerHost time.Duration
}

// Validate returns an error if any duration is non-positive.
func (o Options) Validate() error {
	if o.Dial <= 0 {
		return fmt.Errorf("timeout: Dial must be positive, got %v", o.Dial)
	}
	if o.Exec <= 0 {
		return fmt.Errorf("timeout: Exec must be positive, got %v", o.Exec)
	}
	if o.PerHost <= 0 {
		return fmt.Errorf("timeout: PerHost must be positive, got %v", o.PerHost)
	}
	return nil
}

// ForHost returns a context that is cancelled after the PerHost deadline
// relative to the supplied parent context.
func (o Options) ForHost(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, o.PerHost)
}

// ForExec returns a context that is cancelled after the Exec deadline
// relative to the supplied parent context.
func (o Options) ForExec(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, o.Exec)
}

// ForDial returns a context that is cancelled after the Dial deadline
// relative to the supplied parent context.
func (o Options) ForDial(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, o.Dial)
}
