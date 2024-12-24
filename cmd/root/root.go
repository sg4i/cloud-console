package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{ // 将 rootCmd 改为 RootCmd
		Use:   "cloudconsole",
		Short: "Cloud Console CLI",
		Long:  `A CLI application for managing cloud services.`,
	}
)

// 添加这个新函数
func InitCommands() {
	// 定义命令组
	cliGroup := &cobra.Group{
		ID:    "cli",
		Title: "cli Commands:",
	}
	rpcGroup := &cobra.Group{
		ID:    "server",
		Title: "server Commands:",
	}

	// 将主要命令组添加到根命令
	RootCmd.AddGroup(cliGroup)
	RootCmd.AddGroup(rpcGroup)

	// 这里会导入并初始化所有子命令
	initCliCommands(cliGroup)

}

func Execute() {
	InitCommands()                            // 在执行之前调用初始化函数
	if err := RootCmd.Execute(); err != nil { // 修复这里的语法错误
		fmt.Println(err)
		os.Exit(1)
	}
}
