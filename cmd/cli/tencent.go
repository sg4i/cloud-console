package cli

import (
	"fmt"
	"os"

	"github.com/sg4i/cloud-console/config"
	"github.com/sg4i/cloud-console/pkg/console"
	"github.com/spf13/cobra"
)

func NewTencentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tencent",
		Short: "生成腾讯云角色登录 URL",
		Run:   run(),
	}

	AddStringFlag(cmd, "secret-id", "i", "", "腾讯云 Secret ID", false)
	AddStringFlag(cmd, "secret-key", "k", "", "腾讯云 Secret Key", false)
	AddStringFlag(cmd, "token", "t", "", "腾讯云 Token", false)
	AddStringFlag(cmd, "role-arn", "r", "", "腾讯云角色 ARN", false)
	AddStringFlag(cmd, "destination", "d", "https://console.cloud.tencent.com", "登录成功后跳转的 URL", false)
	AddBoolFlag(cmd, "auto-login", "a", true, "自动打开 URL", false)

	// 关闭参数自动排序
	cmd.Flags().SortFlags = false

	return cmd
}

func run() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// 获取命令行参数
		cfg := config.New()
		provider := cfg.GetProvider().GetTencent()

		secretId, _ := cmd.Flags().GetString("secret-id")
		secretKey, _ := cmd.Flags().GetString("secret-key")
		token, _ := cmd.Flags().GetString("token")
		roleArn, _ := cmd.Flags().GetString("role-arn")
		destination, _ := cmd.Flags().GetString("destination")

		// 如果命令行参数为空，尝试从配置文件读取
		if secretId == "" || secretKey == "" {
			// 从配置文件获取凭证信息
			if credential := provider.GetCredential(); credential != nil {
				if secretId == "" {
					secretId = credential.SecretId
				}
				if secretKey == "" {
					secretKey = credential.SecretKey
				}
				if token == "" {
					token = credential.Token
				}
			}
		}

		// 从配置文件获取其他参数
		if roleArn == "" {
			roleArn = provider.GetRoleArn()
		}
		if destination == "" {
			destination = provider.GetDestination()
		}

		// 校验必要参数
		if secretId == "" || secretKey == "" {
			fmt.Println("错误: Secret ID 和 Secret Key 不能为空")
			os.Exit(1)
		}

		// 创建 Console 实例
		opts := &console.Options{
			Mode:     "cli",
			Provider: console.ProviderTencent,
		}

		c, err := console.New(opts)
		if err != nil {
			fmt.Printf("创建 Console 实例失败: %v\n", err)
			os.Exit(1)
		}

		// 创建登录选项
		loginOpts := console.NewLoginOptions(secretId, secretKey, token, roleArn, destination, "")

		// 获取登录 URL
		url, err := c.GetLoginURL(loginOpts)
		if err != nil {
			fmt.Printf("获取登录 URL 失败: %v\n", err)
			os.Exit(1)
		}

		// 输出登录 URL
		fmt.Println(url)
	}
}
