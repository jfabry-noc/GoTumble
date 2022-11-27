package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
)

type AuthConfig struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	Token          string `json:"token"`
	TokenSecret    string `json:"token_secret"`
	Instance       string `json:"instance"`
	Format         string `json:"format"`
}

// getConfigPath determines the path for the configuration file.
func getConfigPath() (string, error) {
	currentUser, err := user.Current()
	return fmt.Sprintf("%v/.config/gotumble.json", currentUser.HomeDir), err
}

// LoadConfig loads an existing configuration file.
func LoadConfig() (AuthConfig, error) {
	var authInfo AuthConfig
	configPath, err := getConfigPath()
	if err != nil {
		return authInfo, err
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return authInfo, err
	}

	err = json.Unmarshal(content, &authInfo)
	return authInfo, err
}

// WriteConfig writes a new configuration file.
func WriteConfig(config AuthConfig) (string, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Sprintf("Failed to determine config file path with error: %v\n", err), err
	}

	jsonText, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Sprintf("Failed to parse configuration to JSON with error: %v\n", err), err
	}

	err = ioutil.WriteFile(configPath, jsonText, 0600)
	if err != nil {
		return fmt.Sprintf("Failed to write file to %v with error: %v\n", configPath, err), err
	}

	return fmt.Sprintf("Successfully wrote the config file to: %v\n", configPath), nil
}
