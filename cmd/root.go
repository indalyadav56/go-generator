package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Root command is the base command for your CLI
var rootCmd = &cobra.Command{
	Use:   "yourcli",
	Short: "Your Project CLI",
	Long:  "A command-line tool for interacting with YourProject services.",
}

// Execute is the main entry point for all Cobra commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// func init() {
// 	// Here you can add any global or persistent flags for all commands
// 	rootCmd.PersistentFlags().StringP("config", "c", "", "Config file (default is $HOME/.yourcli.yaml)")
// }
