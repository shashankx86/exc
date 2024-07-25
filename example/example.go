package example

import (
	"fmt"

	"exc/config"
	"exc/internal/plugin"
)

type ExamplePlugin struct{}

func (p *ExamplePlugin) Name() string {
	return "example"
}

func (p *ExamplePlugin) Execute(action config.Action, variables map[string]string) error {
	fmt.Println("Executing custom example plugin action")
	return nil
}

func init() {
	plugin.RegisterPlugin(&ExamplePlugin{})
}
