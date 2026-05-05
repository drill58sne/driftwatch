// Package plugin provides a lightweight extension mechanism for driftwatch,
// allowing custom check functions to be registered and executed alongside
// built-in checks during a drift scan.
package plugin

import (
	"errors"
	"fmt"
	"sync"

	"github.com/driftwatch/internal/checker"
)

// CheckFn is the signature that all plugin check functions must satisfy.
type CheckFn func(host string) (checker.CheckResult, error)

// Plugin represents a named, versioned extension.
type Plugin struct {
	Name    string
	Version string
	Check   CheckFn
}

// Registry holds registered plugins and is safe for concurrent use.
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]Plugin
}

// New returns an initialised, empty Registry.
func New() *Registry {
	return &Registry{
		plugins: make(map[string]Plugin),
	}
}

// Register adds p to the registry. Returns an error if a plugin with the
// same name is already registered or if required fields are missing.
func (r *Registry) Register(p Plugin) error {
	if p.Name == "" {
		return errors.New("plugin name must not be empty")
	}
	if p.Check == nil {
		return fmt.Errorf("plugin %q: Check function must not be nil", p.Name)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plugins[p.Name]; exists {
		return fmt.Errorf("plugin %q is already registered", p.Name)
	}

	r.plugins[p.Name] = p
	return nil
}

// Get returns the plugin registered under name, or an error if not found.
func (r *Registry) Get(name string) (Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.plugins[name]
	if !ok {
		return Plugin{}, fmt.Errorf("plugin %q not found", name)
	}
	return p, nil
}

// List returns a snapshot of all registered plugin names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}
	return names
}

// Unregister removes the plugin with the given name. No-op if not present.
func (r *Registry) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.plugins, name)
}
