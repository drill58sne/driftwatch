// Package audit provides structured, append-only audit logging for driftwatch.
//
// Events are written as newline-delimited JSON to a configurable file path,
// making them easy to ingest by log aggregators or query with standard tools.
//
// # Writing events
//
//	l, err := audit.New("/var/log/driftwatch/audit.log")
//	if err != nil { ... }
//	l.Record(audit.Event{
//		Kind:    audit.EventScan,
//		Host:    "web-01",
//		Message: "scan complete",
//		Drifted: true,
//	})
//
// # Reading events
//
//	events, err := audit.ReadAll("/var/log/driftwatch/audit.log")
//	filtered := audit.FilterByHost(events, "web-01")
//	audit.SortByTime(filtered)
package audit
