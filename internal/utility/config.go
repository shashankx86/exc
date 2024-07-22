package utility

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "exc/config"
)

// LoadConfig reads and parses the JSON configuration file
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

    var config config.CommandConfig
    err = json.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }

    return &config, nil
}
