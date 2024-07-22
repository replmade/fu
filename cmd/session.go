package cmd

import (
	"fmt"
	"fu/global"
	"fu/internal"
	"log"

	"github.com/spf13/cobra"
)

var expiresIn int64

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Retrieve a session cookie for the current app's user",
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

		idToken, idTokenOk := appConfig["id_token"].(string)
		if !idTokenOk || idToken == "" {
			log.Fatalf("ID token not found. Please sign in using the `signin` command.")
		}

		global.Fa, err = internal.InitializeFirebase(appConfig)
		if err != nil {
			log.Fatalf(err.Error())
		}
		global.Fa.SetIdToken(idToken)

		sessionCookie, err := global.Fa.GetSessionCookie(expiresIn)
		if err != nil {
			log.Fatalf("failed to get session cookie: %v", err)
		}

		fmt.Printf("session cookie: %s\n", sessionCookie)

		cfg[global.AppName].(map[string]interface{})["session"] = sessionCookie
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
