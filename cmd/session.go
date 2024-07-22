package cmd

import (
	"fmt"
	"fu/global"
	"fu/internal"
	"log"

	fb "github.com/replmade/firebase-spells-go/auth"
	"github.com/spf13/cobra"
)

var expiresIn int64

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Retrieve a session cookie for the current app's user",
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
			log.Fatalf("no configuration found for app %s in %s", currentApp, configFilePath)
		}

		apiKey, apiKeyOk := appConfig["api_key"].(string)
		saKeyPath, saKeyPathOk := appConfig["sa_key_path"].(string)
		global.ApiKey = apiKey
		global.SaKeyPath = saKeyPath
		if !apiKeyOk || !saKeyPathOk {
			log.Fatalf("invalid API key or service account key path for app %s", currentApp)
		}

		idToken, idTokenOk := appConfig["id_token"].(string)
		if !idTokenOk || idToken == "" {
			log.Fatalf("ID token not found. Please sign in using the `signin` command.")
		}

		global.Fa = &fb.FirebaseAuth{}
		global.Fa.SetIdToken(idToken)
		if err := global.Fa.Initialize(global.SaKeyPath); err != nil {
			log.Fatalf("failed to initialize Firebase: %v", err)
		}
		global.Fa.SetAPIKey(global.ApiKey)

		sessionCookie, err := global.Fa.GetSessionCookie(expiresIn)
		if err != nil {
			log.Fatalf("failed to get session cookie: %v", err)
		}

		fmt.Printf("session cookie: %s\n", sessionCookie)

		// Save the session token in the config file under the current app
		appConfig["session"] = sessionCookie
		if err := internal.WriteConfig(configFile, cfg); err != nil {
			log.Fatalf("failed to write config to file: %v", err)
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if expiresIn <= 0 {
			expiresIn = 86400
		}
		return nil
	},
}

func init() {
	sessionCmd.Flags().Int64Var(&expiresIn, "expires-in", 86400, "Session cookie expiration time in seconds")
	rootCmd.AddCommand(sessionCmd)
}
