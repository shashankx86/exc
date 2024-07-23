package main

import (
    "log"
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "exc/cmd"
    "exc/internal/utility"
)

// debug flag will be set at compile time
var dev string

func main() {
    // Determine the config path based on the debug flag
    var configPath string
    if dev == "1" {
        // Dev mode: use local config path
        configPath = "example/.exc.config.json"
    } else {
        // Release mode: use home directory config path
        homeDir, err := os.UserHomeDir()
        if err != nil {
            log.Fatalf("Failed to get home directory: %v", err)
        }
        configPath = filepath.Join(homeDir, ".exc.config.json")
    }

    // Load configuration
    config, err := utility.LoadConfig(configPath)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Create the root command
    rootCmd := &cobra.Command{
        Use:   "exc",
        Short: "A dynamically generated CLI tool using user config",
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
