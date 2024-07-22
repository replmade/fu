package cmd

import (
	"fmt"
	"fu/global"
	"fu/internal"
	"log"

	fb "github.com/replmade/firebase-spells-go/auth"
	"github.com/spf13/cobra"
)

var customTokenCmd = &cobra.Command{
	Use:   "custom",
	Short: "Generate a custom token for the current app's user using the session token",
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
			log.Fatalf("no `settings` section found in %s", configFilePath)
		}

		currentApp, ok := configSection["current-app"].(string)
		if !ok || currentApp == "" {
			log.Fatalf("no current app set in %s", configFilePath)
		}

		appConfig, ok := cfg[currentApp].(map[string]interface{})
		if !ok {
			log.Fatalf("no configuration found for app %s in %s", currentApp, configFilePath)
		}

		sessionCookie, sessionOk := appConfig["session"].(string)
		if !sessionOk || sessionCookie == "" {
			log.Fatalf("session token not found. Please get a session token with the command `session`.")
		}

		apiKey, apiKeyOk := appConfig["api_key"].(string)
		saKeyPath, saKeyPathOk := appConfig["sa_key_path"].(string)
		if !apiKeyOk || !saKeyPathOk {
			log.Fatalf("invalid API key or service account key path for app %s", currentApp)
		}

		global.Fa = &fb.FirebaseAuth{}
		if err := global.Fa.Initialize(saKeyPath); err != nil {
			log.Fatalf("failed to initialize Firebase: %v", err)
		}
		global.Fa.SetAPIKey(apiKey)
		global.Fa.SetSessionCookie(sessionCookie)

		authToken, err := global.Fa.AuthUser()
		if err != nil {
			log.Fatalf("failed to authenticate user: %v", err)
		}

		customToken, err := global.Fa.CreateCustomToken(authToken.UID)
		if err != nil {
			log.Fatalf("failed to get custom token: %v", err)
		}
		fmt.Printf("Custom token: %s\n", customToken)

		appConfig["custom_token"] = customToken
		if err := internal.WriteConfig(configFile, cfg); err != nil {
			log.Fatalf("failed to write custom token to file: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(customTokenCmd)
}
