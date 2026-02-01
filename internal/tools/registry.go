package tools

import "sync"

var (
	registry   = make(map[string]Tool)
	registryMu sync.RWMutex
)

// Register adds a tool by name.
func Register(t Tool) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[t.Name()] = t
}

// Get returns a tool by name, or nil.
func Get(name string) Tool {
	registryMu.RLock()
	defer registryMu.RUnlock()
	return registry[name]
}

// All returns all registered tools.
func All() []Tool {
	registryMu.RLock()
	defer registryMu.RUnlock()
	out := make([]Tool, 0, len(registry))
	for _, t := range registry {
		out = append(out, t)
	}
	return out
}
