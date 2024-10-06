package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main2() {
	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "App is a CLI application example",
	}

	var helloCmd = &cobra.Command{
		Use:   "hello",
		Short: "Prints Hello",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from Cobra!")
		},
	}

	rootCmd.AddCommand(helloCmd)
	rootCmd.Execute()
}
