package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
	case "make_http_request":
		return executeHTTPRequest(action, variables)
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
	return nil
}

// executeHTTPRequest performs an HTTP request based on the action configuration
func executeHTTPRequest(action config.Action, variables map[string]string) error {
	method := strings.ToUpper(action.Method)
	url := action.URL

	// Replace placeholders with variables in URL
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		url = strings.ReplaceAll(url, placeholder, value)
	}

	var reqBody []byte
	if action.Body != "" {
		body := action.Body
		// Replace placeholders with variables in Body
		for key, value := range variables {
			placeholder := fmt.Sprintf("{{.%s}}", key)
			body = strings.ReplaceAll(body, placeholder, value)
		}
		reqBody = []byte(body)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	for key, value := range action.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP response body: %v", err)
	}

	logrus.Infof("HTTP response status: %s", resp.Status)
	logrus.Infof("HTTP response body: %s", string(respBody))
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
