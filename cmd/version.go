package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the CLI",
	Long:  `All software has versions. This is the CLI's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("go-generator v0.1.5")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
