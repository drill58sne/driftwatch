package plugin

import (
	"context"
	"fmt"

	"github.com/driftwatch/internal/checker"
)

// RunOptions controls execution behaviour for RunAll.
type RunOptions struct {
	// StopOnError causes RunAll to abort after the first plugin error.
	StopOnError bool
}

// DefaultRunOptions returns conservative defaults.
func DefaultRunOptions() RunOptions {
	return RunOptions{
		StopOnError: false,
	}
}

// PluginResult pairs a plugin name with its check outcome.
type PluginResult struct {
	Plugin string
	Result checker.CheckResult
	Err    error
}

// RunAll executes every registered plugin against host and collects results.
// The provided context is checked between plugin invocations; if cancelled the
// function returns immediately with the results collected so far.
func (r *Registry) RunAll(ctx context.Context, host string, opts RunOptions) []PluginResult {
	names := r.List()
	results := make([]PluginResult, 0, len(names))

	for _, name := range names {
		if err := ctx.Err(); err != nil {
			results = append(results, PluginResult{
				Plugin: name,
				Err:    fmt.Errorf("context cancelled before running plugin %q: %w", name, err),
			})
			break
		}

		p, err := r.Get(name)
		if err != nil {
			results = append(results, PluginResult{Plugin: name, Err: err})
			if opts.StopOnError {
				break
			}
			continue
		}

		res, err := p.Check(host)
		results = append(results, PluginResult{Plugin: name, Result: res, Err: err})

		if err != nil && opts.StopOnError {
			break
		}
	}

	return results
}
