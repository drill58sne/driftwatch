// Package tagger provides utilities for tagging and grouping check results
// by arbitrary labels such as environment, team, or service.
package tagger

import (
	"fmt"
	"sort"
	"strings"

	"github.com/driftwatch/internal/checker"
)

// Tag represents a key-value label attached to a result.
type Tag struct {
	Key   string
	Value string
}

// String returns the canonical "key=value" representation of a Tag.
func (t Tag) String() string {
	return fmt.Sprintf("%s=%s", t.Key, t.Value)
}

// ParseTag parses a "key=value" string into a Tag.
// It returns an error if the format is invalid.
func ParseTag(s string) (Tag, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Tag{}, fmt.Errorf("tagger: invalid tag format %q, expected key=value", s)
	}
	return Tag{Key: parts[0], Value: parts[1]}, nil
}

// Group partitions results by the value of the given tag key.
// Results that do not carry the tag are placed under the key "(untagged)".
func Group(results []checker.CheckResult, tagKey string) map[string][]checker.CheckResult {
	groups := make(map[string][]checker.CheckResult)
	const untagged = "(untagged)"

	for _, r := range results {
		matched := false
		for _, tag := range r.Tags {
			if tag == tagKey || strings.HasPrefix(tag, tagKey+"=") {
				val := strings.TrimPrefix(tag, tagKey+"=")
				groups[val] = append(groups[val], r)
				matched = true
				break
			}
		}
		if !matched {
			groups[untagged] = append(groups[untagged], r)
		}
	}
	return groups
}

// Keys returns a sorted list of unique tag keys present across all results.
func Keys(results []checker.CheckResult) []string {
	seen := make(map[string]struct{})
	for _, r := range results {
		for _, raw := range r.Tags {
			if k, _, ok := strings.Cut(raw, "="); ok {
				seen[k] = struct{}{}
			}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
