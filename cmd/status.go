package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd represents the "status" command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the current status",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Current status: All systems operational.")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
