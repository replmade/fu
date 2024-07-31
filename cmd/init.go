package cmd

import (
	"fmt"
	"log"

	"fu/global"
	"fu/internal"
	"github.com/replmade/firebase-spells-go/auth"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the Firebase project",
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath, configFile, err := internal.EnsureConfigDirAndFile()
		if err != nil {
			log.Fatalf(err.Error())
		}
		defer configFile.Close()

		cfg, err := internal.ReadConfig(configFile)
		if err != nil {
			log.Fatalf(err.Error())
		}

		cfg[global.AppName] = map[string]string{
			"api_key":     global.ApiKey,
			"sa_key_path": global.SaKeyPath,
		}
		internal.UpdateCurrentApp(cfg, global.AppName)

		if err := internal.WriteConfig(configFile, cfg); err != nil {
			log.Fatalf(err.Error())
		}

		fmt.Printf("Configuration written to %s\n", configFilePath)

		global.Fa = &fb.FirebaseAuth{}
		if err := global.Fa.Initialize(global.SaKeyPath); err != nil {
			log.Fatalf("Failed to initialize Firebase: %v", err)
		}
		global.Fa.SetAPIKey(global.ApiKey)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if global.ApiKey == "" {
			return fmt.Errorf("--api-key is required")
		}
		if global.SaKeyPath == "" {
			return fmt.Errorf("--sa-key-path is required")
		}
		if global.AppName == "" {
			return fmt.Errorf("--app-name is required")
		}
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&global.AppName, "app-name", "", "Firebase app name")
	initCmd.Flags().StringVar(&global.ApiKey, "api-key", "", "Firebase app API Key")
	initCmd.Flags().StringVar(&global.SaKeyPath, "sa-key-path", "", "Path to the service account key JSON file")
	initCmd.MarkFlagRequired("app-name")
	initCmd.MarkFlagRequired("api-key")
	initCmd.MarkFlagRequired("sa-key-path")
	rootCmd.AddCommand(initCmd)
}
