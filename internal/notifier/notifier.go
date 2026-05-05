// Package notifier provides webhook-based drift notifications.
package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/driftwatch/internal/differ"
)

// Config holds configuration for the webhook notifier.
type Config struct {
	WebhookURL string
	Timeout    time.Duration
	Headers    map[string]string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Timeout: 10 * time.Second,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

// Payload is the JSON body sent to the webhook.
type Payload struct {
	Timestamp   time.Time        `json:"timestamp"`
	TotalHosts  int              `json:"total_hosts"`
	DriftedHosts int             `json:"drifted_hosts"`
	Summary     differ.Summary   `json:"summary"`
}

// Notifier sends drift summaries to a configured webhook endpoint.
type Notifier struct {
	cfg    Config
	client *http.Client
}

// New creates a Notifier with the given Config.
func New(cfg Config) *Notifier {
	return &Notifier{
		cfg: cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

// Notify sends a drift summary to the configured webhook URL.
// It returns an error if the request fails or the server responds with a non-2xx status.
func (n *Notifier) Notify(summary differ.Summary) error {
	if n.cfg.WebhookURL == "" {
		return fmt.Errorf("notifier: webhook URL is not configured")
	}

	payload := Payload{
		Timestamp:    time.Now().UTC(),
		TotalHosts:   summary.TotalHosts,
		DriftedHosts: summary.DriftedHosts,
		Summary:      summary,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notifier: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, n.cfg.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notifier: build request: %w", err)
	}

	for k, v := range n.cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("notifier: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notifier: unexpected status %d from webhook", resp.StatusCode)
	}

	return nil
}
