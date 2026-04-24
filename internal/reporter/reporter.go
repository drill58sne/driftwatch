package reporter

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// Format defines the output format for drift reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// DriftResult holds the drift check result for a single host and check.
type DriftResult struct {
	Host      string
	CheckName string
	Expected  string
	Actual    string
	Drifted   bool
	Timestamp time.Time
}

// Reporter writes drift results to an output destination.
type Reporter struct {
	format Format
	out    io.Writer
}

// New creates a new Reporter with the given format, writing to stdout by default.
func New(format Format) *Reporter {
	return &Reporter{
		format: format,
		out:    os.Stdout,
	}
}

// NewWithWriter creates a Reporter that writes to the provided writer.
func NewWithWriter(format Format, w io.Writer) *Reporter {
	return &Reporter{format: format, out: w}
}

// Write outputs the drift results in the configured format.
func (r *Reporter) Write(results []DriftResult) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(results)
	default:
		return r.writeText(results)
	}
}

func (r *Reporter) writeText(results []DriftResult) error {
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "HOST\tCHECK\tSTATUS\tEXPECTED\tACTUAL")
	for _, res := range results {
		status := "OK"
		if res.Drifted {
			status = "DRIFT"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			res.Host, res.CheckName, status, res.Expected, res.Actual)
	}
	return w.Flush()
}

func (r *Reporter) writeJSON(results []DriftResult) error {
	fmt.Fprintln(r.out, "[")
	for i, res := range results {
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		fmt.Fprintf(r.out,
			`  {"host":%q,"check":%q,"drifted":%v,"expected":%q,"actual":%q,"timestamp":%q}%s\n`,
			res.Host, res.CheckName, res.Drifted, res.Expected, res.Actual,
			res.Timestamp.Format(time.RFC3339), comma)
	}
	fmt.Fprintln(r.out, "]")
	return nil
}
