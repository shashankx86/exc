package utility

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"exc/config"

	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

const (
	profileDir     = ".exc/profiles"
	activeProfile  = ".exc/active_profile"
	defaultProfile = "default"
)

// LoadConfig reads and parses the JSON configuration file and validates it against the schema
func LoadConfig(filePath string) (*config.CommandConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := ValidateConfig(data); err != nil {
		return nil, err
	}

	var cfg config.CommandConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ValidateConfig validates the JSON configuration against the schema
func ValidateConfig(configData []byte) error {
	schemaLoader := gojsonschema.NewStringLoader(Schema)
	documentLoader := gojsonschema.NewBytesLoader(configData)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			logrus.Errorf("Validation error: %s", desc)
		}
		return fmt.Errorf("configuration validation failed")
	}

	return nil
}

func GetActiveProfilePath(dev string) string {
	if dev == "1" {
		return "example/.exc.config.json"
	}
	activeProfilePath := filepath.Join(profileDir, getActiveProfile())
	return activeProfilePath
}

func getActiveProfile() string {
	data, err := ioutil.ReadFile(activeProfile)
	if err != nil {
		return defaultProfile
	}
	return string(data)
}

func ListProfiles() ([]string, error) {
	files, err := ioutil.ReadDir(profileDir)
	if err != nil {
		return nil, err
	}
	profiles := make([]string, 0, len(files))
	for _, file := range files {
		profiles = append(profiles, file.Name())
	}
	return profiles, nil
}

func AddProfile(name, configPath string) error {
	profilePath := filepath.Join(profileDir, name)
	if _, err := os.Stat(profilePath); err == nil {
		return fmt.Errorf("profile %s already exists", name)
	}
	return copyFile(configPath, profilePath)
}

func SwitchProfile(name string) error {
	profilePath := filepath.Join(profileDir, name)
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return fmt.Errorf("profile %s does not exist", name)
	}
	return ioutil.WriteFile(activeProfile, []byte(name), 0644)
}

func DeleteProfile(name string) error {
	profilePath := filepath.Join(profileDir, name)
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return fmt.Errorf("profile %s does not exist", name)
	}
	return os.Remove(profilePath)
}

func copyFile(src, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, data, 0644)
}

// Schema represents the JSON schema for validating the configuration
const Schema = `
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "properties": {
        "commands": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "id": { "type": "string" },
                    "description": { "type": "string" },
                    "actions": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "type": { "type": "string" },
                                "message": { "type": "string" },
                                "variable_name": { "type": "string" },
                                "value": { "type": "string" },
                                "url": { "type": "string" },
                                "method": { "type": "string" },
                                "headers": {
                                    "type": "object",
                                    "additionalProperties": { "type": "string" }
                                },
                                "body": { "type": "string" },
                                "condition": { "type": "string" },
                                "trueActions": { "type": "array" },
                                "falseActions": { "type": "array" },
                                "onError": { "type": "string" },
                                "retry": {
                                    "type": "object",
                                    "properties": {
                                        "count": { "type": "integer" },
                                        "interval": { "type": "integer" }
                                    }
                                },
                                "timeout": { "type": "integer" }
                            },
                            "required": ["type"]
                        }
                    },
                    "subcommands": {
                        "type": "array",
                        "items": {
                            "$ref": "#/properties/commands/items"
                        }
                    }
                },
                "required": ["id", "description", "actions"]
            }
        }
    },
    "required": ["commands"]
}
`

func init() {
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		logrus.Fatalf("Failed to create profile directory: %v", err)
	}
}
