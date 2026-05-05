// Package audit provides structured audit logging for drift detection events.
// Each significant action (scan, alert, baseline update) is recorded with
// a timestamp, host, and outcome for later review.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventKind classifies the type of audit event.
type EventKind string

const (
	EventScan     EventKind = "scan"
	EventAlert    EventKind = "alert"
	EventBaseline EventKind = "baseline_update"
	EventError    EventKind = "error"
)

// Event represents a single auditable action.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Kind      EventKind `json:"kind"`
	Host      string    `json:"host"`
	Message   string    `json:"message"`
	Drifted   bool      `json:"drifted"`
}

// Logger writes audit events as newline-delimited JSON.
type Logger struct {
	w io.Writer
}

// New returns a Logger that appends events to the file at path.
// The file is created if it does not exist.
func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("audit: open %s: %w", path, err)
	}
	return &Logger{w: f}, nil
}

// NewWithWriter returns a Logger that writes to w. Useful for testing.
func NewWithWriter(w io.Writer) *Logger {
	return &Logger{w: w}
}

// Record encodes and writes a single audit event.
func (l *Logger) Record(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}
