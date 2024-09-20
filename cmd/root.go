package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cloud-console",
	Short: "Cloud Console is a tool for managing cloud resources",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
	Long: `A CLI application for managing cloud services.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
