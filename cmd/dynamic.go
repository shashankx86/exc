package cmd

import (
    "fmt"

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
                for key, value := range cmdConfig.Actions {
                    if key == "message" {
                        logrus.Debugf("Executing action: %s", value)
                        fmt.Println(value)
                    }
                }
            },
        }

        // Add the new command to the root command
        rootCmd.AddCommand(cmd)
    }
}
