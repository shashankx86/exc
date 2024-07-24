package cmd

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"

    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "exc/config"
)

// GenerateDynamicCommands generates commands based on the configuration
func GenerateDynamicCommands(rootCmd *cobra.Command, config *config.CommandConfig) {
    for _, cmdConfig := range config.Commands {
        logrus.Debugf("Generating command: %s", cmdConfig.ID)
        cmd := createCommand(cmdConfig)
        rootCmd.AddCommand(cmd)
    }
}

// createCommand creates a new Cobra command from the command configuration
func createCommand(cmdConfig config.Command) *cobra.Command {
    return &cobra.Command{
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
}

// executeAction executes a single action
func executeAction(action config.Action, variables map[string]string) error {
    switch action.Type {
    case "print":
        return executePrintAction(action, variables)
    case "set_variable":
        return executeSetVariableAction(action, variables)
    case "make_http_request":
        return executeHTTPRequestAction(action, variables)
    default:
        return fmt.Errorf("unknown action type: %s", action.Type)
    }
}

// executePrintAction handles the print action
func executePrintAction(action config.Action, variables map[string]string) error {
    msg := replacePlaceholders(action.Message, variables)
    fmt.Println(msg)
    return nil
}

// executeSetVariableAction handles the set_variable action
func executeSetVariableAction(action config.Action, variables map[string]string) error {
    variables[action.VariableName] = action.Value
    logrus.Infof("Set variable %s to %s", action.VariableName, action.Value)
    return nil
}

// executeHTTPRequestAction performs an HTTP request based on the action configuration
func executeHTTPRequestAction(action config.Action, variables map[string]string) error {
    method, url, reqBody, err := prepareHTTPRequest(action, variables)
    if err != nil {
        return err
    }

    req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
    if err != nil {
        return fmt.Errorf("failed to create HTTP request: %v", err)
    }

    setRequestHeaders(req, action.Headers)
    resp, err := performHTTPRequest(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read HTTP response body: %v", err)
    }

    logrus.Infof("HTTP response status: %s", resp.Status)
    logrus.Infof("HTTP response body: %s", string(respBody))

    if action.ResponseVar != "" {
        variables[action.ResponseVar] = string(respBody)
        logrus.Infof("Stored HTTP response in variable: %s", action.ResponseVar)
    }

    return nil
}

// prepareHTTPRequest prepares the HTTP request parameters
func prepareHTTPRequest(action config.Action, variables map[string]string) (string, string, []byte, error) {
    method := strings.ToUpper(action.Method)
    url := replacePlaceholders(action.URL, variables)
    var reqBody []byte
    if action.Body != "" {
        body := replacePlaceholders(action.Body, variables)
        reqBody = []byte(body)
    }
    return method, url, reqBody, nil
}

// setRequestHeaders sets the headers for the HTTP request
func setRequestHeaders(req *http.Request, headers map[string]string) {
    for key, value := range headers {
        req.Header.Set(key, value)
    }
}

// performHTTPRequest executes the HTTP request
func performHTTPRequest(req *http.Request) (*http.Response, error) {
    client := &http.Client{}
    return client.Do(req)
}

// replacePlaceholders replaces placeholders in a string with variable values
func replacePlaceholders(input string, variables map[string]string) string {
    for key, value := range variables {
        placeholder := fmt.Sprintf("{{.%s}}", key)
        input = strings.ReplaceAll(input, placeholder, value)
    }
    return input
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
