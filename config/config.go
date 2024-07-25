package config

const CliVersion = "0.0.1"

type CommandConfig struct {
    Commands []Command `json:"commands"`
}

type Command struct {
    ID          string    `json:"id"`
    Aliases     []string  `json:"aliases,omitempty"`
    Description string    `json:"description"`
    Actions     []Action  `json:"actions"`
    Subcommands []Command `json:"subcommands,omitempty"`
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
    Condition    string            `json:"condition,omitempty"`
    TrueActions  []Action          `json:"trueActions,omitempty"`
    FalseActions []Action          `json:"falseActions,omitempty"`
    OnError      string            `json:"onError,omitempty"`
    Retry        Retry             `json:"retry,omitempty"`
    Timeout      int               `json:"timeout,omitempty"`
    LoopCount    int               `json:"loopCount,omitempty"`
    LoopActions  []Action          `json:"loopActions,omitempty"`
    ResponseVar  string            `json:"responseVar,omitempty"`
}

type Retry struct {
    Count    int `json:"count"`
    Interval int `json:"interval"`
}
