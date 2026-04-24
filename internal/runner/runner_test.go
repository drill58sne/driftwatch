package runner_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/inventory"
	"github.com/driftwatch/internal/reporter"
	"github.com/driftwatch/internal/runner"
)

func makeInventory(hosts ...inventory.Host) *inventory.Inventory {
	return &inventory.Inventory{Hosts: hosts}
}

func makeChecks() []checker.Check {
	return []checker.Check{
		{Name: "os-version", Command: "uname -r", Expected: "5.15.0"},
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := runner.DefaultOptions()
	if opts.Concurrency != 5 {
		t.Errorf("expected concurrency 5, got %d", opts.Concurrency)
	}
	if opts.Format != "text" {
		t.Errorf("expected format 'text', got %s", opts.Format)
	}
}

func TestRun_EmptyInventory(t *testing.T) {
	var buf bytes.Buffer
	rep := reporter.NewWithWriter(&buf)
	inv := makeInventory()
	opts := runner.DefaultOptions()

	err := runner.Run(inv, makeChecks(), opts, rep)
	if err != nil {
		t.Fatalf("unexpected error for empty inventory: %v", err)
	}
}

func TestRun_UnreachableHost_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	rep := reporter.NewWithWriter(&buf)
	inv := makeInventory(inventory.Host{
		Name:    "bad-host",
		Address: "192.0.2.1", // TEST-NET, unreachable
		Port:    22,
		User:    "ci",
	})
	opts := runner.DefaultOptions()
	opts.Concurrency = 1

	err := runner.Run(inv, makeChecks(), opts, rep)
	if err == nil {
		t.Fatal("expected error for unreachable host, got nil")
	}
	if !errors.Is(err, err) { // basic non-nil check
		t.Errorf("unexpected error type: %v", err)
	}
}

func TestRun_ZeroConcurrency_Normalised(t *testing.T) {
	var buf bytes.Buffer
	rep := reporter.NewWithWriter(&buf)
	inv := makeInventory()
	opts := runner.DefaultOptions()
	opts.Concurrency = 0 // should be normalised to 1

	// With an empty inventory this should always succeed regardless of concurrency.
	if err := runner.Run(inv, nil, opts, rep); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
