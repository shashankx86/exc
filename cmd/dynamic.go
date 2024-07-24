package cmd

import (
    "fmt"
    "strings"

    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "exc/config"
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
                variables := make(map[string]string)
                for _, action := range cmdConfig.Actions {
                    if err := executeAction(action, variables); err != nil {
                        handleActionError(action, err)
                        if action.OnError == "stop" {
                            break
                        }
                    }
                }
            },
        }

        // Add the new command to the root command
        rootCmd.AddCommand(cmd)
    }
}

// executeAction executes a single action
func executeAction(action config.Action, variables map[string]string) error {
    switch action.Type {
    case "print":
        msg := action.Message
        // Replace placeholders with variables
        for key, value := range variables {
            placeholder := fmt.Sprintf("{{.%s}}", key)
            msg = strings.ReplaceAll(msg, placeholder, value)
        }
        fmt.Println(msg)
    case "set_variable":
        variables[action.VariableName] = action.Value
        logrus.Infof("Set variable %s to %s", action.VariableName, action.Value)
    default:
        return fmt.Errorf("unknown action type: %s", action.Type)
    }
    return nil
}

// handleActionError handles errors based on the action's onError field
func handleActionError(action config.Action, err error) {
    switch action.OnError {
    case "log":
        logrus.Errorf("Error executing action %s: %v", action.Type, err)
    case "stop":
        logrus.Errorf("Error executing action %s: %v. Stopping execution.", action.Type, err)
    default:
        logrus.Errorf("Error executing action %s: %v", action.Type, err)
    }
}
