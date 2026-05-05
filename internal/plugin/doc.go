// Package plugin provides a thread-safe registry for extending driftwatch
// with custom check functions.
//
// Plugins are registered by name and version, and expose a CheckFn that
// accepts a host string and returns a checker.CheckResult. The Registry
// is safe for concurrent reads and writes.
//
// Basic usage:
//
//	reg := plugin.New()
//
//	err := reg.Register(plugin.Plugin{
//		Name:    "kernel-version",
//		Version: "1.0.0",
//		Check: func(host string) (checker.CheckResult, error) {
//			// ... SSH and run `uname -r` ...
//			return checker.CheckResult{Name: "kernel-version", Output: "6.1.0"}, nil
//		},
//	})
//
//	results := reg.RunAll(ctx, "web-01", plugin.DefaultRunOptions())
package plugin
