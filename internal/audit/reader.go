package audit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// ReadAll reads all audit events from the file at path.
func ReadAll(path string) ([]Event, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: open %s: %w", path, err)
	}
	defer f.Close()
	return decode(f)
}

// FilterByHost returns only events matching the given host.
func FilterByHost(events []Event, host string) []Event {
	out := make([]Event, 0, len(events))
	for _, e := range events {
		if e.Host == host {
			out = append(out, e)
		}
	}
	return out
}

// FilterByKind returns only events matching the given kind.
func FilterByKind(events []Event, kind EventKind) []Event {
	out := make([]Event, 0, len(events))
	for _, e := range events {
		if e.Kind == kind {
			out = append(out, e)
		}
	}
	return out
}

// SortByTime sorts events in descending order (newest first).
func SortByTime(events []Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.After(events[j].Timestamp)
	})
}

func decode(r io.Reader) ([]Event, error) {
	var events []Event
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var e Event
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: decode line: %w", err)
		}
		events = append(events, e)
	}
	return events, sc.Err()
}
