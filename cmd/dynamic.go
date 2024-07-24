package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"exc/config"
	"exc/internal/utility"
)

func GenerateDynamicCommands(rootCmd *cobra.Command, config *config.CommandConfig) {
	for _, cmdConfig := range config.Commands {
		logrus.Debugf("Generating command: %s", cmdConfig.ID)
		cmd := createCommand(cmdConfig)
		rootCmd.AddCommand(cmd)
	}
}

func createCommand(cmdConfig config.Command) *cobra.Command {
	return &cobra.Command{
		Use:   cmdConfig.ID,
		Short: cmdConfig.Description,
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
}
