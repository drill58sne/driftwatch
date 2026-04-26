// Package alert implements threshold-based alerting for driftwatch.
//
// An Alerter evaluates a differ.Summary against configurable warn and error
// thresholds, returning an Alert when the number of drifted checks meets or
// exceeds a threshold. Alerts can be emitted to any io.Writer, defaulting to
// stdout.
//
// Example usage:
//
//	cfg := alert.Config{WarnThreshold: 1, ErrorThreshold: 5}
//	a := alert.New(cfg)
//	if al := a.Evaluate(summary); al != nil {
//	    a.Emit(al)
//	}
package alert
