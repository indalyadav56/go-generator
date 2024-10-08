/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("project called")

		// Fetch the name flag value
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name = "World"
		}
		fmt.Printf("Project Name, %s!\n", name)
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)
	projectCmd.Flags().StringP("name", "n", "", "The name to greet")
}
