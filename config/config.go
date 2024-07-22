package config

const CliVersion = "0.0.1"

// CommandConfig represents the structure of the JSON configuration
type CommandConfig struct {
    Commands []Command `json:"commands"`
}

// Command represents a single command configuration
type Command struct {
    ID          string            `json:"id"`
    Description string            `json:"description"`
    Actions     map[string]string `json:"actions"`
}
