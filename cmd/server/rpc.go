package server

import (
	"path/filepath"

	"github.com/sg4i/cloud-console/cmd/common"
	"github.com/sg4i/cloud-console/config"

	"github.com/sg4i/cloud-console/pkg/console"
	"github.com/spf13/cobra"
)

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "启动 RPC 服务器和 HTTP 网关",
		Long:  `启动 RPC 服务器和 HTTP 网关,提供与命令行相同的功能`,
		RunE:  run(),
	}

	// 设置默认配置文件路径
	defaultConfig := filepath.Join(".", "config.yml")

	common.AddStringFlag(cmd, "config", "c", defaultConfig, "配置文件路径", false)
	common.AddStringFlag(cmd, "rpc-address", "", ":50050", "RPC 服务器地址和端口", false)
	common.AddStringFlag(cmd, "http-address", "", ":50080", "HTTP 网关地址和端口", false)
	common.AddBoolFlag(cmd, "no-http", "", false, "不启动 HTTP 网关", false)
	common.AddStringFlag(cmd, "token", "t", "", "认证令牌", false)

	// 绑定命令行参数到变量
	cmd.Flags().SortFlags = false

	return cmd
}

func run() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		configFile, _ := cmd.Flags().GetString("config")
		rpcAddress, _ := cmd.Flags().GetString("rpc-address")
		httpAddress, _ := cmd.Flags().GetString("http-address")
		noHTTP, _ := cmd.Flags().GetBool("token")
		token, _ := cmd.Flags().GetString("token")

		if noHTTP {
			httpAddress = ""
		}

		if token == "" {
			cfg := config.New(configFile)
			if rpcCfg := cfg.GetRPC(); rpcCfg != nil {
				token = rpcCfg.GetToken()
			}
		}

		return console.StartRPCServer(rpcAddress, httpAddress, token)
	}
}
