package rollup

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Writer renders a rollup Summary to an output stream.
type Writer struct {
	w io.Writer
}

// NewWriter returns a Writer that writes to w.
// If w is nil, os.Stdout is used.
func NewWriter(w io.Writer) *Writer {
	if w == nil {
		w = os.Stdout
	}
	return &Writer{w: w}
}

// WriteText writes a human-readable summary table.
func (wr *Writer) WriteText(s Summary) error {
	fmt.Fprintf(wr.w, "Rollup by %-8s  %s\n", s.GroupBy, s.RolledAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(wr.w, "%-20s %6s %7s %5s\n", "KEY", "TOTAL", "DRIFTED", "CLEAN")
	fmt.Fprintln(wr.w, "--------------------------------------------")
	for _, e := range s.Entries {
		status := "ok"
		if e.HasDrift() {
			status = "DRIFT"
		}
		fmt.Fprintf(wr.w, "%-20s %6d %7d %5d  [%s]\n",
			e.Key, e.Total, e.Drifted, e.Clean, status)
	}
	return nil
}

// WriteJSON writes the summary as a JSON document.
func (wr *Writer) WriteJSON(s Summary) error {
	enc := json.NewEncoder(wr.w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}
