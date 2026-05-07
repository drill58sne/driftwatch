package sampler_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/sampler"
)

func makeResults(output string) []checker.Result {
	return []checker.Result{
		{Name: "check1", Host: "host1", Output: output, Drift: output != "ok"},
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := sampler.DefaultOptions()
	if opts.WindowSize != 10 {
		t.Errorf("expected WindowSize 10, got %d", opts.WindowSize)
	}
	if opts.MaxAge != 30*time.Minute {
		t.Errorf("expected MaxAge 30m, got %v", opts.MaxAge)
	}
}

func TestNew_NotNil(t *testing.T) {
	s := sampler.New(sampler.DefaultOptions())
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
}

func TestNew_ClampsBelowOne(t *testing.T) {
	s := sampler.New(sampler.Options{WindowSize: 0, MaxAge: time.Minute})
	for i := 0; i < 5; i++ {
		s.Record("h", makeResults("ok"))
	}
	samples := s.Get("h")
	if len(samples) != 1 {
		t.Errorf("expected 1 sample (clamped window), got %d", len(samples))
	}
}

func TestRecord_And_Get(t *testing.T) {
	s := sampler.New(sampler.Options{WindowSize: 5, MaxAge: time.Minute})
	s.Record("host1", makeResults("ok"))
	s.Record("host1", makeResults("drift"))

	samples := s.Get("host1")
	if len(samples) != 2 {
		t.Fatalf("expected 2 samples, got %d", len(samples))
	}
	if samples[1].Results[0].Output != "drift" {
		t.Errorf("unexpected output: %s", samples[1].Results[0].Output)
	}
}

func TestGet_UnknownHost_ReturnsEmpty(t *testing.T) {
	s := sampler.New(sampler.DefaultOptions())
	samples := s.Get("nonexistent")
	if len(samples) != 0 {
		t.Errorf("expected 0 samples, got %d", len(samples))
	}
}

func TestRecord_EvictsExcessSamples(t *testing.T) {
	s := sampler.New(sampler.Options{WindowSize: 3, MaxAge: time.Minute})
	for i := 0; i < 7; i++ {
		s.Record("h", makeResults("ok"))
	}
	if got := len(s.Get("h")); got != 3 {
		t.Errorf("expected 3 samples after eviction, got %d", got)
	}
}

func TestRecord_EvictsStaleByAge(t *testing.T) {
	s := sampler.New(sampler.Options{WindowSize: 10, MaxAge: 50 * time.Millisecond})
	s.Record("h", makeResults("ok"))
	time.Sleep(80 * time.Millisecond)
	samples := s.Get("h")
	if len(samples) != 0 {
		t.Errorf("expected stale sample to be evicted, got %d", len(samples))
	}
}

func TestHosts_ReturnsRecordedHosts(t *testing.T) {
	s := sampler.New(sampler.DefaultOptions())
	s.Record("alpha", makeResults("ok"))
	s.Record("beta", makeResults("ok"))

	hosts := s.Hosts()
	if len(hosts) != 2 {
		t.Errorf("expected 2 hosts, got %d", len(hosts))
	}
}
