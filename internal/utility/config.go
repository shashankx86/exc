package utility

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"exc/config"

	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

// LoadConfig reads and parses the JSON configuration file and validates it against the schema
func LoadConfig(filePath string) (*config.CommandConfig, error) {
	// Read the config file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Validate the configuration against the schema
	if err := ValidateConfig(data); err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into CommandConfig struct
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
                    }
                },
                "required": ["id", "description", "actions"]
            }
        }
    },
    "required": ["commands"]
}
`
