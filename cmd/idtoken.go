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
		_, appConfig, err := internal.LoadAppConfig()
		if err != nil {
			log.Fatalf(err.Error())
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
