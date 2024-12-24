package root

import (
	"github.com/sg4i/cloud-console/cmd/cli"
	"github.com/spf13/cobra"
)

func initCliCommands(group *cobra.Group) {
	cliCmd := cli.NewTencentCmd()
	cliCmd.GroupID = group.ID
	RootCmd.AddCommand(cliCmd)
}
