package config

const CliVersion = "0.0.1"

// CommandConfig represents the overall configuration for commands
type CommandConfig struct {
	Commands []Command `json:"commands"`
}

// Command represents a single command in the configuration
type Command struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Actions     []Action `json:"actions"`
}

type Action struct {
	Type         string            `json:"type"`
	Message      string            `json:"message,omitempty"`
	VariableName string            `json:"variable_name,omitempty"`
	Value        string            `json:"value,omitempty"`
	URL          string            `json:"url,omitempty"`
	Method       string            `json:"method,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	Body         string            `json:"body,omitempty"`
	ResponseVar  string            `json:"response_var,omitempty"`
	OnError      string            `json:"onError,omitempty"`
}
