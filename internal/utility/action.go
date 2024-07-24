package utility

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"exc/config"
)

func ExecuteAction(action config.Action, variables map[string]string) error {
	switch action.Type {
	case "print":
		return executePrintAction(action, variables)
	case "set_variable":
		return executeSetVariableAction(action, variables)
	case "make_http_request":
		return executeHTTPRequestAction(action, variables)
	case "condition":
		return executeConditionalAction(action, variables)
	case "loop":
		return executeLoopAction(action, variables)
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

func executePrintAction(action config.Action, variables map[string]string) error {
	msg := ReplacePlaceholders(action.Message, variables)
	fmt.Println(msg)
	return nil
}

func executeSetVariableAction(action config.Action, variables map[string]string) error {
	variables[action.VariableName] = action.Value
	logrus.Infof("Set variable %s to %s", action.VariableName, action.Value)
	return nil
}

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

func executeConditionalAction(action config.Action, variables map[string]string) error {
	condition := ReplacePlaceholders(action.Condition, variables)
	if EvalCondition(condition) {
		for _, trueAction := range action.TrueActions {
			if err := ExecuteAction(trueAction, variables); err != nil {
				return err
			}
		}
	} else {
		for _, falseAction := range action.FalseActions {
			if err := ExecuteAction(falseAction, variables); err != nil {
				return err
			}
		}
	}
	return nil
}

func executeLoopAction(action config.Action, variables map[string]string) error {
	for i := 0; i < action.LoopCount; i++ {
		for _, loopAction := range action.LoopActions {
			if err := ExecuteAction(loopAction, variables); err != nil {
				HandleActionError(loopAction, err)
				if loopAction.OnError == "stop" {
					return err
				}
			}
		}
	}
	return nil
}

func prepareHTTPRequest(action config.Action, variables map[string]string) (string, string, []byte, error) {
	method := strings.ToUpper(action.Method)
	url := ReplacePlaceholders(action.URL, variables)
	var reqBody []byte
	if action.Body != "" {
		body := ReplacePlaceholders(action.Body, variables)
		reqBody = []byte(body)
	}
	return method, url, reqBody, nil
}

func setRequestHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func performHTTPRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

func ReplacePlaceholders(input string, variables map[string]string) string {
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		input = strings.ReplaceAll(input, placeholder, value)
	}
	input = replaceEnvPlaceholders(input)
	return input
}

func replaceEnvPlaceholders(input string) string {
	startIdx := strings.Index(input, "{{.Env.")
	for startIdx != -1 {
		endIdx := strings.Index(input[startIdx:], "}}")
		if endIdx == -1 {
			break
		}
		endIdx += startIdx + 2
		placeholder := input[startIdx:endIdx]
		envVarName := strings.TrimPrefix(strings.TrimSuffix(placeholder, "}}"), "{{.Env.")
		envVarValue := os.Getenv(envVarName)
		input = strings.ReplaceAll(input, placeholder, envVarValue)
		startIdx = strings.Index(input, "{{.Env.")
	}
	return input
}

func EvalCondition(condition string) bool {
	return strings.ToLower(condition) == "true"
}

func HandleActionError(action config.Action, err error) {
	switch action.OnError {
	case "log":
		logrus.Errorf("Error executing action %s: %v", action.Type, err)
	case "stop":
		logrus.Errorf("Error executing action %s: %v. Stopping execution.", action.Type, err)
	default:
		logrus.Errorf("Error executing action %s: %v", action.Type, err)
	}
}
