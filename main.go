package main

import (
	"os"
	"path/filepath"

	"exc/cmd"
	"exc/internal/utility"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// dev flag will be set at compile time
var dev string
var verbose bool

func main() {
	// Get the CLI tool's name
	cliName := filepath.Base(os.Args[0])

	// Create the root command
	rootCmd := &cobra.Command{
		Use:   cliName,
		Short: "A dynamically generated CLI tool using user config",
	}

	// Add the --verbose flag
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Set up logging in a PersistentPreRun function
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if dev == "1" {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("Running in DEV mode")
		} else if verbose {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("Verbose logging enabled")
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}
	}

	// Determine the config path based on the dev flag and active profile
	profilePath := utility.GetActiveProfilePath(dev)

	logrus.Debugf("Using config path: %s", profilePath)

	// Load configuration
	config, err := utility.LoadConfig(profilePath)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Generate commands dynamically
	cmd.GenerateDynamicCommands(rootCmd, config)

	// Add version command
	rootCmd.AddCommand(cmd.NewVersionCommand())

	// Add profile management commands
	rootCmd.AddCommand(cmd.NewProfileCommand())

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Error executing command: %v", err)
	}
}
