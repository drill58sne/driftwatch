// Package redactor provides utilities for masking sensitive values
// in check results before they are written to output or stored.
package redactor

import (
	"regexp"
	"strings"

	"github.com/driftwatch/internal/checker"
)

// DefaultPatterns contains common patterns for sensitive data.
var DefaultPatterns = []string{
	`(?i)password`,
	`(?i)secret`,
	`(?i)token`,
	`(?i)api[_-]?key`,
	`(?i)private[_-]?key`,
}

const masked = "[REDACTED]"

// Redactor masks sensitive values in check results.
type Redactor struct {
	patterns []*regexp.Regexp
}

// New creates a Redactor using the provided glob patterns.
// If patterns is nil or empty, DefaultPatterns are used.
func New(patterns []string) (*Redactor, error) {
	if len(patterns) == 0 {
		patterns = DefaultPatterns
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return &Redactor{patterns: compiled}, nil
}

// Apply returns a copy of results with sensitive Output values masked.
func (r *Redactor) Apply(results []checker.CheckResult) []checker.CheckResult {
	out := make([]checker.CheckResult, len(results))
	for i, res := range results {
		out[i] = res
		if r.isSensitive(res.Name) {
			out[i].Output = masked
			out[i].Expected = masked
		}
	}
	return out
}

// isSensitive reports whether name matches any registered pattern.
func (r *Redactor) isSensitive(name string) bool {
	norm := strings.TrimSpace(name)
	for _, re := range r.patterns {
		if re.MatchString(norm) {
			return true
		}
	}
	return false
}
