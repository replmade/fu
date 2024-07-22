package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "fu",
	Short: "Firebase authentication utility",
	Long:  "Authenticate a Firebase user using email and password and retrieve tokens",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(signinCmd)
	rootCmd.AddCommand(idTokenCmd)
	rootCmd.AddCommand(sessionCmd)
	rootCmd.AddCommand(uidCmd)
}
