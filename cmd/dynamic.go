package cmd

import (
	"fmt"

	"exc/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// VariableStore holds variables for the CLI
var VariableStore = make(map[string]string)

// GenerateDynamicCommands generates commands based on the configuration
func GenerateDynamicCommands(rootCmd *cobra.Command, config *config.CommandConfig) {
	for _, cmdConfig := range config.Commands {
		logrus.Debugf("Generating command: %s", cmdConfig.ID)

		// Create a new command
		cmd := &cobra.Command{
			Use:   cmdConfig.ID,
			Short: cmdConfig.Description,
			Run: func(cmd *cobra.Command, args []string) {
				executeActions(cmdConfig.Actions)
			},
		}

		// Add the new command to the root command
		rootCmd.AddCommand(cmd)
	}
}

// executeActions executes the actions defined in the configuration
func executeActions(actions []config.Action) {
	for _, action := range actions {
		switch action.Type {
		case "print":
			handlePrintAction(action)
		case "set_variable":
			handleSetVariableAction(action)
		// Add more cases here for different action types
		default:
			logrus.Warnf("Unknown action type: %s", action.Type)
		}
	}
}

// handlePrintAction handles the print action type
func handlePrintAction(action config.Action) {
	if action.Message != "" {
		fmt.Println(action.Message)
	} else {
		logrus.Warn("Print action has no message to display")
	}
}

// handleSetVariableAction handles the set_variable action type
func handleSetVariableAction(action config.Action) {
	if action.VariableName != "" && action.Value != "" {
		VariableStore[action.VariableName] = action.Value
		logrus.Infof("Set variable %s to %s", action.VariableName, action.Value)
	} else {
		logrus.Warn("set_variable action is missing variable_name or value")
	}
}

// GetVariable retrieves the value of a variable from the store
func GetVariable(name string) string {
	return VariableStore[name]
}

// SetVariable sets the value of a variable in the store
func SetVariable(name, value string) {
	VariableStore[name] = value
}
