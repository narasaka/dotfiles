package plugins

import (
	"fmt"
	"sync"
)

// Registry holds all registered provider plugins.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

var globalRegistry = &Registry{
	providers: make(map[string]Provider),
}

// Register adds a provider plugin to the global registry.
func Register(p Provider) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.providers[p.Name()] = p
}

// Get returns a provider by name from the global registry.
func Get(name string) (Provider, error) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	p, ok := globalRegistry.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %q not registered", name)
	}
	return p, nil
}

// List returns all registered provider names and display names.
func List() []ProviderInfo {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	var result []ProviderInfo
	for _, p := range globalRegistry.providers {
		result = append(result, ProviderInfo{
			Name:        p.Name(),
			DisplayName: p.DisplayName(),
		})
	}
	return result
}

type ProviderInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}
