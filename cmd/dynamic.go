package cmd

import (
	"fmt"

	"exc/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
