package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

// getBaseConfigPath determines the base path for the configuration file.
func getBaseConfigPath() (string, error) {
	currentUser, err := user.Current()
	return currentUser.HomeDir, err
}

// getFullConfigPath determines the path for the configuration file.
func getFullConfigPath() (string, error) {
	currentUser, err := user.Current()
	return fmt.Sprintf("%v/.config/gotumble/config.json", currentUser.HomeDir), err
}

// buildConfigPath ensures each stop in the config path exist and returns the full value.
func buildConfigPath(basePath string, nextFolders []string) (string, error) {
	err := addToFilesystem(basePath, "directory")
	if err != nil {
		return "", err
	}

	for _, v := range nextFolders {
		basePath = fmt.Sprintf("%v/%v", basePath, v)

		err = addToFilesystem(basePath, "directory")
		if err != nil {
			return "", err
		}
	}

	basePath = fmt.Sprintf("%v/config.json", basePath)
	return basePath, nil
}

// LoadConfig loads an existing configuration file.
func LoadConfig() (AuthConfig, error) {
	var authInfo AuthConfig
	configPath, err := getFullConfigPath()
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

// addToFilesystem adds a new file or directory if it doesn't already exist.
func addToFilesystem(filePath string, objectType string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		if objectType == "file" {
			_, err = os.Create(filePath)
		} else {
			err = os.Mkdir(filePath, 0755)
		}
	}

	return err
}

// WriteConfig writes a new configuration file.
func WriteConfig(config AuthConfig) (string, error) {
	configPath, err := getBaseConfigPath()
	if err != nil {
		return fmt.Sprintf("Failed to determine config file path with error: %v\n", err), err
	}

	configPath, err = buildConfigPath(configPath, []string{".config", "gotumble"})
	if err != nil {
		return fmt.Sprintf("Failed to build directory path with error: %v\n", err), err
	}

	jsonText, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Sprintf("Failed to parse configuration to JSON with error: %v\n", err), err
	}

	err = addToFilesystem(configPath, "file")
	if err != nil {
		return fmt.Sprintf("Failed to create the file at %v with error: %v\n", configPath, err), err
	}

	err = ioutil.WriteFile(configPath, jsonText, 0600)
	if err != nil {
		return fmt.Sprintf("Failed to write to file at %v with error: %v\n", configPath, err), err
	}

	return fmt.Sprintf("Successfully wrote the config file to: %v\n", configPath), nil
}
