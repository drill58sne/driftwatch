package audit_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/audit"
)

func makeEvent(kind audit.EventKind, host string, drifted bool) audit.Event {
	return audit.Event{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Host:      host,
		Message:   "test event",
		Drifted:   drifted,
	}
}

func TestRecord_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewWithWriter(&buf)
	err := l.Record(makeEvent(audit.EventScan, "web-01", false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "\"kind\":\"scan\"") {
		t.Errorf("expected JSON to contain kind=scan, got: %s", buf.String())
	}
}

func TestRecord_SetsTimestamp(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewWithWriter(&buf)
	e := audit.Event{Kind: audit.EventAlert, Host: "db-01"}
	if err := l.Record(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "timestamp") {
		t.Errorf("expected timestamp in output, got: %s", buf.String())
	}
}

func TestRecord_MultipleEvents(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewWithWriter(&buf)
	for i := 0; i < 3; i++ {
		if err := l.Record(makeEvent(audit.EventScan, "host", false)); err != nil {
			t.Fatalf("record %d: %v", i, err)
		}
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestFilterByHost(t *testing.T) {
	events := []audit.Event{
		makeEvent(audit.EventScan, "web-01", false),
		makeEvent(audit.EventScan, "db-01", true),
		makeEvent(audit.EventAlert, "web-01", true),
	}
	got := audit.FilterByHost(events, "web-01")
	if len(got) != 2 {
		t.Errorf("expected 2 events, got %d", len(got))
	}
}

func TestFilterByKind(t *testing.T) {
	events := []audit.Event{
		makeEvent(audit.EventScan, "web-01", false),
		makeEvent(audit.EventAlert, "web-01", true),
		makeEvent(audit.EventBaseline, "db-01", false),
	}
	got := audit.FilterByKind(events, audit.EventAlert)
	if len(got) != 1 {
		t.Errorf("expected 1 event, got %d", len(got))
	}
}

func TestSortByTime(t *testing.T) {
	now := time.Now()
	events := []audit.Event{
		{Timestamp: now.Add(-2 * time.Minute), Host: "old"},
		{Timestamp: now, Host: "new"},
		{Timestamp: now.Add(-1 * time.Minute), Host: "mid"},
	}
	audit.SortByTime(events)
	if events[0].Host != "new" {
		t.Errorf("expected newest first, got %s", events[0].Host)
	}
}
