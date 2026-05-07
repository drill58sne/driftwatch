// Package sampler provides periodic result sampling for drift trend analysis.
// It collects checker results at configurable intervals and retains a fixed
// window of samples per host for lightweight in-memory trend tracking.
package sampler

import (
	"sync"
	"time"

	"github.com/driftwatch/internal/checker"
)

// DefaultOptions returns a Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		WindowSize: 10,
		MaxAge:     30 * time.Minute,
	}
}

// Options controls sampling behaviour.
type Options struct {
	// WindowSize is the maximum number of samples retained per host.
	WindowSize int
	// MaxAge is the maximum age of a sample before it is evicted.
	MaxAge time.Duration
}

// Sample holds a single point-in-time snapshot of results for a host.
type Sample struct {
	CapturedAt time.Time
	Results    []checker.Result
}

// Sampler records and retrieves result samples per host.
type Sampler struct {
	mu      sync.Mutex
	opts    Options
	windows map[string][]Sample
}

// New returns a new Sampler with the given options.
// If opts.WindowSize is less than 1 it is clamped to 1.
func New(opts Options) *Sampler {
	if opts.WindowSize < 1 {
		opts.WindowSize = 1
	}
	if opts.MaxAge <= 0 {
		opts.MaxAge = DefaultOptions().MaxAge
	}
	return &Sampler{
		opts:    opts,
		windows: make(map[string][]Sample),
	}
}

// Record appends a sample for the given host, evicting stale or excess entries.
func (s *Sampler) Record(host string, results []checker.Result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	sample := Sample{CapturedAt: now, Results: results}

	win := s.windows[host]
	win = append(win, sample)
	win = s.evict(win, now)
	s.windows[host] = win
}

// Get returns all current samples for the given host.
func (s *Sampler) Get(host string) []Sample {
	s.mu.Lock()
	defer s.mu.Unlock()

	win := s.evict(s.windows[host], time.Now())
	s.windows[host] = win

	out := make([]Sample, len(win))
	copy(out, win)
	return out
}

// Hosts returns the list of hosts that have at least one live sample.
func (s *Sampler) Hosts() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	hosts := make([]string, 0, len(s.windows))
	for h := range s.windows {
		hosts = append(hosts, h)
	}
	return hosts
}

// evict removes samples that are too old or exceed the window size.
func (s *Sampler) evict(win []Sample, now time.Time) []Sample {
	cutoff := now.Add(-s.opts.MaxAge)
	filtered := win[:0]
	for _, sm := range win {
		if sm.CapturedAt.After(cutoff) {
			filtered = append(filtered, sm)
		}
	}
	if len(filtered) > s.opts.WindowSize {
		filtered = filtered[len(filtered)-s.opts.WindowSize:]
	}
	return filtered
}
