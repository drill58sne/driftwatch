// Package runner ties together inventory loading, SSH connectivity, check
// execution, and report generation into a single orchestration layer.
//
// Basic usage:
//
//	inv, _ := inventory.Load("inventory.yaml")
//	checks := []checker.Check{
//		{Name: "kernel", Command: "uname -r", Expected: "5.15.0"},
//	}
//	var buf bytes.Buffer
//	rep := reporter.NewWithWriter(&buf)
//	err := runner.Run(inv, checks, runner.DefaultOptions(), rep)
//
// Concurrency controls how many hosts are checked simultaneously. The default
// value is 5. Setting it to 0 is treated as 1 (sequential execution).
//
// Format selects the report output format; valid values are "text" and "json".
package runner
