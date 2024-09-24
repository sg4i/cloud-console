package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	autoLogin bool
	rootCmd   = &cobra.Command{
		Use:   "cloud-console",
		Short: "Cloud Console CLI",
		Long:  `A CLI application for managing cloud services.`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&autoLogin, "auto-login", true, "自动打开 URL")
}
