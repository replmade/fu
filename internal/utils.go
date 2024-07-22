// /internal/utils.go

package internal

import (
	"fmt"
	"fu/global"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	fb "github.com/replmade/firebase-spells-go/auth"
)

// Ensures that the config directory and file exist.
func EnsureConfigDirAndFile() (string, *os.File, error) {
	configDir := filepath.Join(os.Getenv("HOME"), ".fu")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	configFilePath := filepath.Join(configDir, "config.toml")
	configFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create/open config file: %v", err)
	}

	return configFilePath, configFile, nil
}

// Reads the configuration from the file.
func ReadConfig(configFile *os.File) (map[string]interface{}, error) {
	config := make(map[string]interface{})
	if stat, _ := configFile.Stat(); stat.Size() > 0 {
		if err := toml.NewDecoder(configFile).Decode(&config); err != nil {
			return nil, fmt.Errorf("failed to read config from file: %v", err)
		}
	}

	return config, nil
}

// Writes the configuration to the file.
func WriteConfig(configFile *os.File, config map[string]interface{}) error {
	configFile.Seek(0, 0)
	configFile.Truncate(0)
	if err := toml.NewEncoder(configFile).Encode(config); err != nil {
		return fmt.Errorf("failed to write config to file: %v", err)
	}
	return nil
}

// Updates the name of the currently selected app in the configuration.
func UpdateCurrentApp(config map[string]interface{}, appName string) {
	if _, exists := config["settings"]; !exists {
		config["settings"] = map[string]interface{}{}
	}
	configSection := config["settings"].(map[string]interface{})
	configSection["current-app"] = appName
}

// Loads the app configuration from the config file.
func LoadAppConfig() (map[string]interface{}, map[string]interface{}, error) {
	_, configFile, err := EnsureConfigDirAndFile()
	if err != nil {
		return nil, nil, err
	}
	defer configFile.Close()

	cfg, err := ReadConfig(configFile)
	if err != nil {
		return nil, nil, err
	}

	settingsSection, ok := cfg["settings"].(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("no settings section found in config file")
	}

	currentApp, ok := settingsSection["current-app"].(string)
	if !ok || currentApp == "" {
		return nil, nil, fmt.Errorf("no current app set in config file")
	}

	global.AppName = currentApp

	appConfig, ok := cfg[currentApp].(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("no configuration found for app %s in config file", currentApp)
	}

	return cfg, appConfig, nil
}

// Initializes the Firebase authentication with the provided app config.
func InitializeFirebase(appConfig map[string]interface{}) (*fb.FirebaseAuth, error) {
	apiKey, apiKeyOk := appConfig["api_key"].(string)
	saKeyPath, saKeyPathOk := appConfig["sa_key_path"].(string)
	if !apiKeyOk || !saKeyPathOk {
		return nil, fmt.Errorf("invalid API key or service account key path")
	}

	fa := &fb.FirebaseAuth{}
	if err := fa.Initialize(saKeyPath); err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase: %v", err)
	}
	fa.SetAPIKey(apiKey)
	return fa, nil
}
