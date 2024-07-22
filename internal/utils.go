package internal

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"os"
	"path/filepath"
)

// Ensure config directory and file exist
func EnsureConfigDirAndFile() (string, *os.File, error) {
	configDir := filepath.Join(os.Getenv("HOME"), ".fu")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", nil, fmt.Errorf("Failed to create config directory: %v", err)
	}

	configFilePath := filepath.Join(configDir, "config.toml")
	configFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to create/open config file: %v", err)
	}

	return configFilePath, configFile, nil
}

// Read the configuration from the file
func ReadConfig(configFile *os.File) (map[string]interface{}, error) {
	config := make(map[string]interface{})
	if stat, _ := configFile.Stat(); stat.Size() > 0 {
		if err := toml.NewDecoder(configFile).Decode(&config); err != nil {
			return nil, fmt.Errorf("failed to read config from file: %v", err)
		}
	}

	return config, nil
}

// Write the configuration to the file
func WriteConfig(configFile *os.File, config map[string]interface{}) error {
	configFile.Seek(0, 0)
	configFile.Truncate(0)
	if err := toml.NewEncoder(configFile).Encode(config); err != nil {
		return fmt.Errorf("failed to write config to file: %v", err)
	}
	return nil
}
