// Package output provides formatting utilities for rendering drift results
// in multiple output formats suitable for terminal display or machine consumption.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/driftwatch/internal/differ"
)

// Format represents the output format for drift results.
type Format string

const (
	FormatText    Format = "text"
	FormatJSON    Format = "json"
	FormatCompact Format = "compact"
)

// Formatter writes drift summaries to an output destination.
type Formatter struct {
	format Format
	w      io.Writer
}

// New creates a Formatter writing to stdout with the given format.
func New(format Format) *Formatter {
	return NewWithWriter(format, os.Stdout)
}

// NewWithWriter creates a Formatter writing to the provided writer.
func NewWithWriter(format Format, w io.Writer) *Formatter {
	if w == nil {
		w = os.Stdout
	}
	return &Formatter{format: format, w: w}
}

// Write renders the drift summary using the configured format.
func (f *Formatter) Write(summary differ.Summary) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(summary)
	case FormatCompact:
		return f.writeCompact(summary)
	default:
		return f.writeText(summary)
	}
}

func (f *Formatter) writeText(s differ.Summary) error {
	status := "clean"
	if s.HasDrift() {
		status = "DRIFT DETECTED"
	}
	_, err := fmt.Fprintf(f.w, "Status : %s\nTotal  : %d\nDrifted: %d\nClean  : %d\n",
		status, s.Total, s.Drifted, s.Clean)
	return err
}

func (f *Formatter) writeCompact(s differ.Summary) error {
	parts := []string{
		fmt.Sprintf("total=%d", s.Total),
		fmt.Sprintf("drifted=%d", s.Drifted),
		fmt.Sprintf("clean=%d", s.Clean),
	}
	_, err := fmt.Fprintln(f.w, strings.Join(parts, " "))
	return err
}

func (f *Formatter) writeJSON(s differ.Summary) error {
	enc := json.NewEncoder(f.w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}
