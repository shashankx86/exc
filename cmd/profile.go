package cmd

import (
	"fmt"

	"exc/internal/utility"
	"github.com/spf13/cobra"
)

func NewProfileCommand() *cobra.Command {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage configuration profiles",
	}

	profileCmd.AddCommand(newProfileListCommand())
	profileCmd.AddCommand(newProfileAddCommand())
	profileCmd.AddCommand(newProfileSwitchCommand())
	profileCmd.AddCommand(newProfileDeleteCommand())

	return profileCmd
}

func newProfileListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		Run: func(cmd *cobra.Command, args []string) {
			profiles, err := utility.ListProfiles()
			if err != nil {
				fmt.Printf("Failed to list profiles: %v\n", err)
				return
			}
			for _, profile := range profiles {
				fmt.Println(profile)
			}
		},
	}
}

func newProfileAddCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "add [name] [config-path]",
		Short: "Add a new profile",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			configPath := args[1]
			if err := utility.AddProfile(name, configPath); err != nil {
				fmt.Printf("Failed to add profile: %v\n", err)
				return
			}
			fmt.Println("Profile added successfully")
		},
	}
}

func newProfileSwitchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "switch [name]",
		Short: "Switch to a different profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if err := utility.SwitchProfile(name); err != nil {
				fmt.Printf("Failed to switch profile: %v\n", err)
				return
			}
			fmt.Println("Profile switched successfully")
		},
	}
}

func newProfileDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if err := utility.DeleteProfile(name); err != nil {
				fmt.Printf("Failed to delete profile: %v\n", err)
				return
			}
			fmt.Println("Profile deleted successfully")
		},
	}
}
