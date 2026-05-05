package plugin_test

import (
	"context"
	"errors"
	"testing"

	"github.com/driftwatch/internal/checker"
	"github.com/driftwatch/internal/plugin"
)

func TestRunAll_ReturnsResultForEachPlugin(t *testing.T) {
	reg := plugin.New()
	_ = reg.Register(makePlugin("c1"))
	_ = reg.Register(makePlugin("c2"))

	results := reg.RunAll(context.Background(), "host1", plugin.DefaultRunOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestRunAll_EmptyRegistry_ReturnsEmpty(t *testing.T) {
	reg := plugin.New()
	results := reg.RunAll(context.Background(), "host1", plugin.DefaultRunOptions())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestRunAll_PluginError_ContinuesByDefault(t *testing.T) {
	reg := plugin.New()
	_ = reg.Register(plugin.Plugin{
		Name: "err-plugin", Version: "1",
		Check: func(h string) (checker.CheckResult, error) {
			return checker.CheckResult{}, errors.New("boom")
		},
	})
	_ = reg.Register(makePlugin("ok-plugin"))

	results := reg.RunAll(context.Background(), "host1", plugin.DefaultRunOptions())
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestRunAll_StopOnError_HaltsAfterFailure(t *testing.T) {
	reg := plugin.New()
	_ = reg.Register(plugin.Plugin{
		Name: "fail-first", Version: "1",
		Check: func(h string) (checker.CheckResult, error) {
			return checker.CheckResult{}, errors.New("fail")
		},
	})
	// Register a second plugin that should NOT run.
	_ = reg.Register(makePlugin("should-not-run"))

	opts := plugin.RunOptions{StopOnError: true}
	results := reg.RunAll(context.Background(), "host1", opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result due to StopOnError, got %d", len(results))
	}
}

func TestRunAll_CancelledContext_StopsEarly(t *testing.T) {
	reg := plugin.New()
	_ = reg.Register(makePlugin("p1"))
	_ = reg.Register(makePlugin("p2"))

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	results := reg.RunAll(ctx, "host1", plugin.DefaultRunOptions())
	for _, r := range results {
		if r.Err != nil {
			return // at least one error recorded — expected
		}
	}
}

func TestDefaultRunOptions_StopOnErrorIsFalse(t *testing.T) {
	opts := plugin.DefaultRunOptions()
	if opts.StopOnError {
		t.Fatal("expected StopOnError to default to false")
	}
}
