package main

import (
    "log"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "exc/config"
    "exc/cmd"
    "exc/internal/utility"
)

func main() {
    // Create the root command
    rootCmd := &cobra.Command{
        Use:     "exc",
        Short:   "A dynamically generated CLI tool using user config",
        Version: config.CliVersion, // Set CLI version from config
    }

    // Load configuration from the home directory
    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatalf("Failed to get home directory: %v", err)
    }
    configPath := filepath.Join(homeDir, ".exc.config.json")
    config, err := utility.LoadConfig(configPath)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Generate commands dynamically
    cmd.GenerateDynamicCommands(rootCmd, config)

    // Add version command
    rootCmd.AddCommand(cmd.NewVersionCommand())

    // Execute the root command
    if err := rootCmd.Execute(); err != nil {
        log.Fatalf("Error executing command: %v", err)
    }
}
