package cmd

import (
	"fmt"
	"fu/global"
	"fu/internal"
	"log"

	"github.com/spf13/cobra"
)

var customTokenCmd = &cobra.Command{
	Use:   "custom",
	Short: "Generate a custom token for the current app's user using the session token",
	Run: func(cmd *cobra.Command, args []string) {
		_, configFile, err := internal.EnsureConfigDirAndFile()
		if err != nil {
			log.Fatalf(err.Error())
		}
		defer configFile.Close()

		cfg, appConfig, err := internal.LoadAppConfig()
		if err != nil {
			log.Fatalf(err.Error())
		}

		sessionCookie, sessionOk := appConfig["session"].(string)
		if !sessionOk || sessionCookie == "" {
			log.Fatalf("session token not found. Please get a session token with the command `session`.")
		}

		global.Fa, err = internal.InitializeFirebase(appConfig)
		if err != nil {
			log.Fatalf(err.Error())
		}
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

		cfg[global.AppName].(map[string]interface{})["custom_token"] = customToken
		if err := internal.WriteConfig(configFile, cfg); err != nil {
			log.Fatalf("failed to write custom token to file: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(customTokenCmd)
}
