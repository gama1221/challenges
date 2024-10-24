package cmd

import (
	"todo-api/cmd/bootstrap"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the server",
	Long:  `Starts the server`,
	Run: func(cmd *cobra.Command, args []string) {
		ReadConfigFromRootCmd()
		bootstrap.Setup()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
