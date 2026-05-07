package summarizer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// Writer renders a Summary to an output stream.
type Writer struct {
	w      io.Writer
	format string // "text" or "json"
}

// NewWriter creates a Writer. If w is nil, os.Stdout is used.
// format must be "text" or "json"; defaults to "text".
func NewWriter(w io.Writer, format string) *Writer {
	if w == nil {
		w = os.Stdout
	}
	if format != "json" {
		format = "text"
	}
	return &Writer{w: w, format: format}
}

// Write outputs the summary in the configured format.
func (wr *Writer) Write(s Summary) error {
	if wr.format == "json" {
		return wr.writeJSON(s)
	}
	return wr.writeText(s)
}

func (wr *Writer) writeText(s Summary) error {
	_, err := fmt.Fprintf(wr.w, "Summary: total=%d drifted=%d clean=%d errored=%d\n",
		s.Total, s.Drifted, s.Clean, s.Errored)
	if err != nil {
		return err
	}

	sorted := make([]HostSummary, len(s.Hosts))
	copy(sorted, s.Hosts)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Host < sorted[j].Host })

	for _, hs := range sorted {
		status := "clean"
		if hs.Drifted > 0 {
			status = "DRIFTED"
		}
		_, err = fmt.Fprintf(wr.w, "  %-20s %s  (drift=%.0f%% checks=%d)\n",
			hs.Host, status, hs.DriftRate*100, hs.Total)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wr *Writer) writeJSON(s Summary) error {
	return json.NewEncoder(wr.w).Encode(s)
}
