package cmd

import (
	"fmt"
	"fu/internal"
	"github.com/spf13/cobra"
	"log"
)

var idTokenCmd = &cobra.Command{
	Use:   "id-token",
	Short: "Get the ID token for the current app's user",
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

		configSection, ok := cfg["settings"].(map[string]interface{})
		if !ok {
			log.Fatalf("no current app is set in %s", configFilePath)
		}

		currentApp, ok := configSection["current-app"].(string)
		if !ok || currentApp == "" {
			log.Fatalf("no current app set in %s", configFilePath)
		}

		appConfig, ok := cfg[currentApp].(map[string]interface{})
		if !ok {
			log.Fatalf("no configuration found for app %s in %s", currentApp, configFilePath)
		}

		idToken, idTokenOk := appConfig["id_token"].(string)
		if !idTokenOk || idToken == "" {
			fmt.Println("ID token not found. Please sign in.")
		} else {
			fmt.Printf("ID Token: %s\n", idToken)
		}
	},
}

func init() {
	rootCmd.AddCommand(idTokenCmd)
}
