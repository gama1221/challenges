package cmd

import (
	"os"
	"todo-api/internal/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todoapp",
	Short: "Todo app for managing user tasks",
	Long:  ``,
	// Run: func(cmd *cobra.Command, args []string) { },
}

var configFile string

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "file", "f", "config.yml", "required config file to read configuration")
}

func ReadConfigFromRootCmd() {
	_, err := config.ReadConfig(configFile)
	if err != nil {
		panic(err)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
