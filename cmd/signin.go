package cmd

import (
	"fmt"
	"fu/global"
	"fu/internal"
	"github.com/spf13/cobra"
	fb "github.com/replmade/firebase-spells-go/auth"
	"log"
)

var (
	userEmail    string
	userPassword string
)

var signinCmd = &cobra.Command{
	Use:   "signin",
	Short: "Sign in a Firebase user with email and password and retrieve an id token",
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

		settingsSection, ok := cfg["settings"].(map[string]interface{})
		if !ok {
			log.Fatalf("no settings section found in %s", configFilePath)
		}

		currentApp, ok := settingsSection["current-app"].(string)
		if !ok || currentApp == "" {
			log.Fatalf("no current app set in %s", configFilePath)
		}

		appConfig, ok := cfg[currentApp].(map[string]interface{})
		if !ok {
			log.Fatalf("no configuration found for app %s", currentApp)
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

		idToken, err := global.Fa.AuthenticateUser(userEmail, userPassword)
		if err != nil {
			log.Fatalf("failed to authenticate user: %v", err)
		}

		fmt.Printf("ID Token retrieved. Use `id-token` command to show the token.\n")
		
		appConfig["id_token"] = idToken
		if err := internal.WriteConfig(configFile, cfg); err != nil {
			log.Fatalf("failed to write id_token to file: %v", err)
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if userEmail == "" {
			return fmt.Errorf("--email is required")
		}

		if userPassword == "" {
			return fmt.Errorf("--password is required")
		}
		return nil
	},
}

func init() {
	signinCmd.Flags().StringVar(&userEmail, "email", "", "User email")
	signinCmd.Flags().StringVar(&userPassword, "password", "", "User password")
	signinCmd.MarkFlagRequired("email")
	signinCmd.MarkFlagRequired("password")
}
