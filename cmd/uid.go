package cmd

import (
	"fmt"
	"fu/global"
	"fu/internal"
	"log"

	"github.com/spf13/cobra"
)

var uidCmd = &cobra.Command{
	Use:   "uid",
	Short: "Retrieve the UID of the authenticated user",
	Run: func(cmd *cobra.Command, args []string) {
		_, configFile, err := internal.EnsureConfigDirAndFile()
		if err != nil {
			log.Fatalf(err.Error())
		}
		defer configFile.Close()

		_, appConfig, err := internal.LoadAppConfig()
		if err != nil {
			log.Fatalf(err.Error())
		}

		sessionCookie, sessionOk := appConfig["session"].(string)
		if !sessionOk || sessionCookie == "" {
			fmt.Println("session token not found. Please get a session token with the command `session`.")
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

		fmt.Printf("UID: %s\n", authToken.UID)
	},
}

func init() {
	rootCmd.AddCommand(uidCmd)
}
