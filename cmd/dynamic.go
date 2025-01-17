package cmd

import (
	"exc/config"
	"exc/internal/utility"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func GenerateDynamicCommands(rootCmd *cobra.Command, config *config.CommandConfig) {
	for _, cmdConfig := range config.Commands {
		logrus.Debugf("Generating command: %s", cmdConfig.ID)
		cmd := createCommand(cmdConfig)
		rootCmd.AddCommand(cmd)
	}
}

func createCommand(cmdConfig config.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:     cmdConfig.ID,
		Aliases: cmdConfig.Aliases,
		Short:   cmdConfig.Description,
		Run: func(cmd *cobra.Command, args []string) {
			variables := make(map[string]string)
			for _, action := range cmdConfig.Actions {
				if err := utility.ExecuteAction(action, variables); err != nil {
					utility.HandleActionError(action, err)
					if action.OnError == "stop" {
						break
					}
				}
			}
		},
	}

	for _, subCmdConfig := range cmdConfig.Subcommands {
		subCmd := createCommand(subCmdConfig)
		cmd.AddCommand(subCmd)
	}

	return cmd
}
