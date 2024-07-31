package cmd

import (
	"fmt"
	"fu/global"
	"fu/internal"
	fb "github.com/replmade/firebase-spells-go/auth"
	"github.com/spf13/cobra"
	"log"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load a configuration for the specified Firebase app",
	Run: func(cmd *cobra.Command, args []string) {
		_, configFile, err := internal.EnsureConfigDirAndFile()
		if err != nil {
			log.Fatalf(err.Error())
		}
		defer configFile.Close()

		cfg, err := internal.ReadConfig(configFile)
		if err != nil {
			log.Fatalf(err.Error())
		}

		if appConfig, ok := cfg[global.AppName].(map[string]interface{}); ok {
			apiKey, apiKeyOk := appConfig["api_key"].(string)
			saKeyPath, saKeyPathOk := appConfig["sa_key_path"].(string)
			if apiKeyOk && saKeyPathOk {
				fmt.Printf("Loaded configuration for app %s\n", global.AppName)

				global.Fa = &fb.FirebaseAuth{}
				if err := global.Fa.Initialize(saKeyPath); err != nil {
					log.Fatalf("Failed to initialize Firebase: %v", err)
				}
				global.Fa.SetAPIKey(apiKey)

				internal.UpdateCurrentApp(cfg, global.AppName)

				if err := internal.WriteConfig(configFile, cfg); err != nil {
					log.Fatalf(err.Error())
				}
			} else {
				log.Printf("App: %s\nCould not find valid API Key or SA Key Path\n", global.AppName)
			}
		} else {
			log.Printf("App: %s not found in configuration\n", global.AppName)
		}
	},
}

func init() {
	loadCmd.Flags().StringVar(&global.AppName, "app-name", "", "Firebase app name")
	loadCmd.MarkFlagRequired("name")
	rootCmd.AddCommand(loadCmd)
}
