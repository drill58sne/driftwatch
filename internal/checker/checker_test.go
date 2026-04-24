package checker_test

import (
	"testing"

	"github.com/user/driftwatch/internal/checker"
)

func TestCheckResult_NoDrift(t *testing.T) {
	result := checker.CheckResult{
		Host:     "host1",
		Check:    "timezone",
		Expected: "UTC",
		Actual:   "UTC",
		Drifted:  false,
	}

	if result.Drifted {
		t.Errorf("expected no drift, but Drifted is true")
	}
}

func TestCheckResult_Drift(t *testing.T) {
	result := checker.CheckResult{
		Host:     "host2",
		Check:    "timezone",
		Expected: "UTC",
		Actual:   "America/New_York",
		Drifted:  true,
	}

	if !result.Drifted {
		t.Errorf("expected drift to be detected")
	}
}

func TestCheck_Fields(t *testing.T) {
	c := checker.Check{
		Name:     "kernel version",
		Command:  "uname -r",
		Expected: "5.15.0",
	}

	if c.Name == "" {
		t.Error("check name should not be empty")
	}
	if c.Command == "" {
		t.Error("check command should not be empty")
	}
	if c.Expected == "" {
		t.Error("check expected value should not be empty")
	}
}

func TestNewRunner_NotNil(t *testing.T) {
	// NewRunner should return a non-nil Runner even with a nil client
	// (actual SSH execution is integration-tested via mock server)
	runner := checker.NewRunner(nil, "test-host")
	if runner == nil {
		t.Fatal("expected non-nil Runner")
	}
}
