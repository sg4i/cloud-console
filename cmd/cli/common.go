package cli

import (
	"github.com/spf13/cobra"
)

func AddStringFlag(cmd *cobra.Command, name, shorthand, defaultValue, usage string, required bool) {
	cmd.Flags().StringP(name, shorthand, defaultValue, usage)
	if required {
		cmd.MarkFlagRequired(name)
	}
}

func AddStringSliceFlag(cmd *cobra.Command, name, shorthand string, defaultValue []string, usage string, required bool) {
	cmd.Flags().StringSliceP(name, shorthand, defaultValue, usage)
	if required {
		cmd.MarkFlagRequired(name)
	}
}

// AddBoolFlag 添加布尔类型的命令行参数
func AddBoolFlag(cmd *cobra.Command, name, shorthand string, defaultValue bool, usage string, required bool) {
	cmd.Flags().BoolP(name, shorthand, defaultValue, usage)
	if required {
		cmd.MarkFlagRequired(name)
	}
}
