package cmd

import (
	"fmt"
	"fu/global"
	"github.com/spf13/cobra"
	"log"
)

var (
	userEmail    string
	userPassword string
)

var signinCmd = &cobra.Command{
	Use:   "signin",
	Short: "Sign in a Firebase user with email and password and retrieve and id token",
	Run: func(cmd *cobra.Command, args []string) {
		if global.Fa == nil {
			log.Fatalf("FirebaseAuth instance is not initialized. Please run the init or load command first.")
		}

		_, err := global.Fa.AuthenticateUser(userEmail, userPassword)
		if err != nil {
			log.Fatalf("Failed to sign in user: %v", err)
		}

		fmt.Printf("ID Token retrieved. Use `id-token` command to show the token.\n")
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
