package server

import (
	"fmt"
	"log"

	"github.com/sg4i/cloud-console/pkg/console"
	"github.com/spf13/cobra"
)

var (
	rpcAddress  string
	httpAddress string
	noHTTP      bool
)

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "启动 RPC 服务器和 HTTP 网关",
		Long:  `启动 RPC 服务器和 HTTP 网关,提供与命令行相同的功能`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("正在启动 RPC 服务器,地址: %s\n", rpcAddress)

			if noHTTP {
				httpAddress = ""
			} else {
				fmt.Printf("正在启动 HTTP 网关,地址: %s\n", httpAddress)
			}

			if err := console.StartRPCServer(rpcAddress, httpAddress, ""); err != nil {
				log.Fatalf("启动服务器失败: %v", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&rpcAddress, "rpc-address", "", ":50050", "RPC 服务器地址和端口")
	cmd.Flags().StringVarP(&httpAddress, "http-address", "", ":50080", "HTTP 网关地址和端口")
	cmd.Flags().BoolVar(&noHTTP, "no-http", false, "不启动 HTTP 网关")

	cmd.Flags().SortFlags = false

	return cmd
}
