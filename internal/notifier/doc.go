// Package notifier provides webhook-based notifications for drift events.
//
// When config drift is detected across remote servers, the notifier can
// dispatch a structured JSON payload to a configured HTTP webhook endpoint.
// This enables integration with external alerting systems such as Slack
// incoming webhooks, PagerDuty event APIs, or custom HTTP listeners.
//
// Basic usage:
//
//	cfg := notifier.DefaultConfig()
//	cfg.WebhookURL = "https://hooks.example.com/drift"
//	n := notifier.New(cfg)
//	if err := n.Notify(summary); err != nil {
//		log.Printf("webhook notification failed: %v", err)
//	}
//
// The payload includes a timestamp, total host count, drifted host count,
// and the full differ.Summary for downstream processing.
package notifier
