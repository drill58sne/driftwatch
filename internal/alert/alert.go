// Package alert provides threshold-based alerting for config drift results.
package alert

import (
	"fmt"
	"io"
	"os"

	"github.com/user/driftwatch/internal/differ"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Config holds alerting thresholds.
type Config struct {
	// WarnThreshold triggers a WARN alert when drifted checks >= value.
	WarnThreshold int
	// ErrorThreshold triggers an ERROR alert when drifted checks >= value.
	ErrorThreshold int
}

// DefaultConfig returns sensible default thresholds.
func DefaultConfig() Config {
	return Config{
		WarnThreshold:  1,
		ErrorThreshold: 5,
	}
}

// Alert represents a triggered alert.
type Alert struct {
	Level   Level
	Message string
	Count   int
}

// Alerter evaluates drift summaries and emits alerts.
type Alerter struct {
	cfg    Config
	writer io.Writer
}

// New creates an Alerter writing to stdout with the given config.
func New(cfg Config) *Alerter {
	return &Alerter{cfg: cfg, writer: os.Stdout}
}

// NewWithWriter creates an Alerter writing to the provided writer.
func NewWithWriter(cfg Config, w io.Writer) *Alerter {
	return &Alerter{cfg: cfg, writer: w}
}

// Evaluate inspects a differ.Summary and returns any triggered Alert, or nil.
func (a *Alerter) Evaluate(s differ.Summary) *Alert {
	count := s.Drifted
	if count <= 0 {
		return nil
	}

	var level Level
	switch {
	case count >= a.cfg.ErrorThreshold:
		level = LevelError
	case count >= a.cfg.WarnThreshold:
		level = LevelWarn
	default:
		return nil
	}

	return &Alert{
		Level:   level,
		Count:   count,
		Message: fmt.Sprintf("%d drifted check(s) detected", count),
	}
}

// Emit writes the alert to the configured writer if non-nil.
func (a *Alerter) Emit(alert *Alert) {
	if alert == nil {
		return
	}
	fmt.Fprintf(a.writer, "[%s] %s\n", alert.Level, alert.Message)
}

// EvaluateAndEmit is a convenience method that evaluates the given summary
// and immediately emits any resulting alert. It returns the Alert that was
// emitted, or nil if no threshold was exceeded.
func (a *Alerter) EvaluateAndEmit(s differ.Summary) *Alert {
	alert := a.Evaluate(s)
	a.Emit(alert)
	return alert
}
