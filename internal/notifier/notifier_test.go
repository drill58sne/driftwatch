package notifier_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/driftwatch/internal/differ"
	"github.com/driftwatch/internal/notifier"
)

func makeSummary() differ.Summary {
	return differ.Summary{
		TotalHosts:   3,
		DriftedHosts: 1,
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := notifier.DefaultConfig()
	if cfg.Timeout != 10*time.Second {
		t.Errorf("expected 10s timeout, got %v", cfg.Timeout)
	}
	if cfg.Headers["Content-Type"] != "application/json" {
		t.Error("expected Content-Type header to be set")
	}
}

func TestNotify_Success(t *testing.T) {
	var received notifier.Payload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("expected Content-Type: application/json")
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := notifier.DefaultConfig()
	cfg.WebhookURL = server.URL
	n := notifier.New(cfg)

	if err := n.Notify(makeSummary()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.TotalHosts != 3 {
		t.Errorf("expected TotalHosts=3, got %d", received.TotalHosts)
	}
	if received.DriftedHosts != 1 {
		t.Errorf("expected DriftedHosts=1, got %d", received.DriftedHosts)
	}
}

func TestNotify_NonOKStatus_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := notifier.DefaultConfig()
	cfg.WebhookURL = server.URL
	n := notifier.New(cfg)

	if err := n.Notify(makeSummary()); err == nil {
		t.Error("expected error for non-2xx status")
	}
}

func TestNotify_EmptyURL_ReturnsError(t *testing.T) {
	cfg := notifier.DefaultConfig()
	n := notifier.New(cfg)
	if err := n.Notify(makeSummary()); err == nil {
		t.Error("expected error when webhook URL is empty")
	}
}

func TestNotify_InvalidURL_ReturnsError(t *testing.T) {
	cfg := notifier.DefaultConfig()
	cfg.WebhookURL = "http://127.0.0.1:0/unreachable"
	n := notifier.New(cfg)
	if err := n.Notify(makeSummary()); err == nil {
		t.Error("expected error for unreachable URL")
	}
}
