// Package dedup provides result deduplication to suppress repeated drift
// alerts for the same host/check combination across consecutive runs.
package dedup

import (
	"sync"
	"time"

	"github.com/driftwatch/internal/checker"
)

// Entry holds the last-seen fingerprint and when it was recorded.
type Entry struct {
	Digest    string
	SeenAt    time.Time
	Suppressed int
}

// Store deduplicates CheckResults by host+name digest.
type Store struct {
	mu      sync.Mutex
	entries map[string]*Entry
	ttl     time.Duration
}

// New returns a Store with the given TTL. After TTL expires the entry is
// evicted and the next identical result will be emitted again.
func New(ttl time.Duration) *Store {
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	return &Store{
		entries: make(map[string]*Entry),
		ttl:     ttl,
	}
}

// Filter returns only the results that have not been seen since the TTL
// window. Results whose digest has changed are always passed through.
func (s *Store) Filter(results []checker.CheckResult) []checker.CheckResult {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	var out []checker.CheckResult
	for _, r := range results {
		key := r.Host + "|" + r.Name
		dig := r.Output

		e, ok := s.entries[key]
		if !ok || now.Sub(e.SeenAt) > s.ttl || e.Digest != dig {
			s.entries[key] = &Entry{Digest: dig, SeenAt: now}
			out = append(out, r)
			continue
		}
		e.Suppressed++
	}
	return out
}

// Stats returns the current entry count and total suppression count.
func (s *Store) Stats() (entries, suppressed int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range s.entries {
		entries++
		suppressed += e.Suppressed
	}
	return
}

// Evict removes all entries whose TTL has elapsed.
func (s *Store) Evict() {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, e := range s.entries {
		if now.Sub(e.SeenAt) > s.ttl {
			delete(s.entries, k)
		}
	}
}
