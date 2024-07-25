package plugin

import (
	"fmt"
	"sync"
)

var (
	mu      sync.Mutex
	plugins = make(map[string]Plugin)
)

// RegisterPlugin registers a new plugin
func RegisterPlugin(p Plugin) {
	mu.Lock()
	defer mu.Unlock()
	plugins[p.Name()] = p
}

// GetPlugin retrieves a plugin by name
func GetPlugin(name string) (Plugin, error) {
	mu.Lock()
	defer mu.Unlock()
	if plugin, ok := plugins[name]; ok {
		return plugin, nil
	}
	return nil, fmt.Errorf("plugin not found: %s", name)
}
