package plugin

import "exc/config"

// Plugin is the interface that custom plugins must implement
type Plugin interface {
	Name() string
	Execute(action config.Action, variables map[string]string) error
}
