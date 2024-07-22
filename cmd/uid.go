package cmd

import (
	"fmt"
	"fu/global"
	"fu/internal"
	"log"

	fb "github.com/replmade/firebase-spells-go/auth"
	"github.com/spf13/cobra"
)

var uidCmd = &cobra.Command{
	Use:   "uid",
	Short: "Retrieve the UID of the authenticated user",
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
			log.Fatalf("no config section found in %s", configFilePath)
		}

		currentApp, ok := configSection["current-app"].(string)
		if !ok || currentApp == "" {
			log.Fatalf("no current app set in %s", configFilePath)
		}

		appConfig, ok := cfg[currentApp].(map[string]interface{})
		if !ok {
			log.Fatalf("No configuration found for app %s in %s", currentApp, configFilePath)
		}

		sessionCookie, sessionOk := appConfig["session"].(string)
		if !sessionOk || sessionCookie == "" {
			fmt.Println("session token not found. Please get a session token with the command `session`.")
		}

		apiKey, apiKeyOk := appConfig["api_key"].(string)
		saKeyPath, saKeyPathOk := appConfig["sa_key_path"].(string)
		if !apiKeyOk || !saKeyPathOk {
			log.Fatalf("Invalid API key or service account key path for app %s", currentApp)
		}

		global.Fa = &fb.FirebaseAuth{}
		if err := global.Fa.Initialize(saKeyPath); err != nil {
			log.Fatalf("Failed to initialize Firebase: %v", err)
		}
		global.Fa.SetAPIKey(apiKey)
		global.Fa.SetSessionCookie(sessionCookie)

		authToken, err := global.Fa.AuthUser()
		if err != nil {
			log.Fatalf("Failed to authenticate user: %v", err)
		}

		fmt.Printf("UID: %s\n", authToken.UID)
	},
}

func init() {
	rootCmd.AddCommand(uidCmd)
}
