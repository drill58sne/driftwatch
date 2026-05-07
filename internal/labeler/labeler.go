// Package labeler provides utilities for attaching and resolving dynamic
// labels to check results based on configurable rules.
package labeler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/driftwatch/internal/checker"
)

// Rule maps a label key/value pair to a set of match conditions.
type Rule struct {
	// Label is the key=value string to attach, e.g. "env=production".
	Label string
	// HostPattern is an optional regex matched against the result host.
	HostPattern string
	// CheckPattern is an optional regex matched against the result check name.
	CheckPattern string
}

// Labeler applies label rules to check results.
type Labeler struct {
	rules []compiledRule
}

type compiledRule struct {
	label       string
	hostRe      *regexp.Regexp
	checkRe     *regexp.Regexp
}

// New creates a Labeler from the provided rules.
// Returns an error if any pattern fails to compile.
func New(rules []Rule) (*Labeler, error) {
	compiled := make([]compiledRule, 0, len(rules))
	for _, r := range rules {
		if strings.TrimSpace(r.Label) == "" {
			return nil, fmt.Errorf("labeler: rule has empty label")
		}
		cr := compiledRule{label: r.Label}
		if r.HostPattern != "" {
			re, err := regexp.Compile(r.HostPattern)
			if err != nil {
				return nil, fmt.Errorf("labeler: invalid host pattern %q: %w", r.HostPattern, err)
			}
			cr.hostRe = re
		}
		if r.CheckPattern != "" {
			re, err := regexp.Compile(r.CheckPattern)
			if err != nil {
				return nil, fmt.Errorf("labeler: invalid check pattern %q: %w", r.CheckPattern, err)
			}
			cr.checkRe = re
		}
		compiled = append(compiled, cr)
	}
	return &Labeler{rules: compiled}, nil
}

// Apply attaches matching labels to each result and returns the annotated slice.
// Labels are appended to the result's existing Tags field.
func (l *Labeler) Apply(results []checker.Result) []checker.Result {
	out := make([]checker.Result, len(results))
	for i, r := range results {
		for _, rule := range l.rules {
			if rule.hostRe != nil && !rule.hostRe.MatchString(r.Host) {
				continue
			}
			if rule.checkRe != nil && !rule.checkRe.MatchString(r.Name) {
				continue
			}
			if !containsTag(r.Tags, rule.label) {
				r.Tags = append(r.Tags, rule.label)
			}
		}
		out[i] = r
	}
	return out
}

func containsTag(tags []string, label string) bool {
	for _, t := range tags {
		if t == label {
			return true
		}
	}
	return false
}
