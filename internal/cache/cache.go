// Package cache provides a lightweight in-memory result cache for driftwatch,
// allowing repeated checks against the same host to reuse recent results
// within a configurable TTL window.
package cache

import (
	"sync"
	"time"

	"github.com/driftwatch/internal/checker"
)

// Entry holds a cached set of check results along with the time they were stored.
type Entry struct {
	Results   []checker.CheckResult
	StoredAt  time.Time
}

// Cache is a thread-safe in-memory store keyed by host address.
type Cache struct {
	mu      sync.RWMutex
	store   map[string]Entry
	ttl     time.Duration
	nowFunc func() time.Time
}

// DefaultTTL is the default time-to-live for cached entries.
const DefaultTTL = 5 * time.Minute

// New returns a Cache with the given TTL. If ttl is zero, DefaultTTL is used.
func New(ttl time.Duration) *Cache {
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	return &Cache{
		store:   make(map[string]Entry),
		ttl:     ttl,
		nowFunc: time.Now,
	}
}

// Set stores results for the given host, overwriting any existing entry.
func (c *Cache) Set(host string, results []checker.CheckResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[host] = Entry{
		Results:  results,
		StoredAt: c.nowFunc(),
	}
}

// Get returns the cached results for host if they exist and have not expired.
// The second return value reports whether a valid entry was found.
func (c *Cache) Get(host string) ([]checker.CheckResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.store[host]
	if !ok {
		return nil, false
	}
	if c.nowFunc().Sub(entry.StoredAt) > c.ttl {
		return nil, false
	}
	return entry.Results, true
}

// Invalidate removes the cached entry for the given host, if any.
func (c *Cache) Invalidate(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, host)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]Entry)
}

// Size returns the current number of entries in the cache, including expired ones.
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}
