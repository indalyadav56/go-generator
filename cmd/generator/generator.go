package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "A CLI tool for managing apps",
	Long:  `A command-line tool that allows managing internal apps.`,
}

// Execute executes the root command, adding subcommands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addAppCmd)
}

func main() {
	Execute()
}

// Define the addApp command
var addAppCmd = &cobra.Command{
	Use:   "addapp [title]",
	Short: "Add a new app with a specific title",
	Long:  `This command adds a new app by creating necessary folders and files dynamically based on the title.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		fmt.Println("title: ", title)
	},
}
