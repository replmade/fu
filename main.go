package main

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/replmade/firebase-spells-go/auth"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var rootCmd = &cobra.Command{
	Use:   "fu",
	Short: "Firebase authentication utility",
	Long:  "Authenticate a Firebase user using email and password and retrieve tokens",
}

var (
	appName   string
	apiKey    string
	saKeyPath string
	fa        *fb.FirebaseAuth
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Firebase project",
	Run: func(cmd *cobra.Command, args []string) {
		configDir := filepath.Join(os.Getenv("HOME"), ".fu")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Fatalf("Failed to create config directory: %v", err)
		}

		configFilePath := filepath.Join(configDir, "config.toml")
		configFile, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalf("Failed to create/open config file: %v", err)
		}
		defer configFile.Close()

		// Read existing config if any
		config := make(map[string]interface{})
		if stat, _ := configFile.Stat(); stat.Size() > 0 {
			if err := toml.NewDecoder(configFile).Decode(&config); err != nil {
				log.Fatalf("Failed to read config from file: %v", err)
			}
		}

		// Write the config to the file
		config[appName] = map[string]string{
			"api_key":     apiKey,
			"sa_key_path": saKeyPath,
		}
		configFile.Seek(0, 0)
		if err := toml.NewEncoder(configFile).Encode(config); err != nil {
			log.Fatalf("Failed to write config to file: %v", err)
		}

		fmt.Printf("Configuration written to %s\n", configFilePath)

		fa = &fb.FirebaseAuth{}
		if err := fa.Initialize(saKeyPath); err != nil {
			log.Fatalf("Failed to initialize Firebase: %v", err)
		}
		fa.SetAPIKey(apiKey)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if apiKey == "" {
			return fmt.Errorf("--api-key is required")
		}
		if saKeyPath == "" {
			return fmt.Errorf("--sa-key-path is required")
		}
		return nil
	},
}

func main() {
	initCmd.Flags().StringVar(&appName, "app-name", "", "Firebase app name")
	initCmd.Flags().StringVar(&apiKey, "api-key", "", "Firebase app API Key")
	initCmd.Flags().StringVar(&saKeyPath, "sa-key-path", "", "Path to the service account key JSON file")
	initCmd.MarkFlagRequired("app-name")
	initCmd.MarkFlagRequired("api-key")
	initCmd.MarkFlagRequired("sa-key-path")

	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
