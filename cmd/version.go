package cmd

import (
    "github.com/spf13/cobra"
    "exc/config"
)

// NewVersionCommand creates a new version command
func NewVersionCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "version",
        Short: "Print the version number of the CLI",
        Run: func(cmd *cobra.Command, args []string) {
            cmd.Println("CLI version:", config.CliVersion)
        },
    }
}
