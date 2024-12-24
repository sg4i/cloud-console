package root

import (
	"github.com/sg4i/cloud-console/cmd/cli"
	"github.com/spf13/cobra"
)

func initCliCommands(group *cobra.Group) {
	tencentCmd := cli.NewTencentCmd()
	tencentCmd.GroupID = group.ID
	RootCmd.AddCommand(tencentCmd)

	alibabaCmd := cli.NewAlibabaCmd()
	alibabaCmd.GroupID = group.ID
	RootCmd.AddCommand(alibabaCmd)

	awsCmd := cli.NewAWSCmd()
	awsCmd.GroupID = group.ID
	RootCmd.AddCommand(awsCmd)
}
