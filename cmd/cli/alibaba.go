package cli

import (
	"fmt"
	"os"

	"github.com/sg4i/cloud-console/config"
	"github.com/sg4i/cloud-console/pkg/console"
	"github.com/spf13/cobra"
)

func NewAlibabaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alibaba",
		Short: "生成阿里云角色登录 URL",
		Run:   runAlibaba(),
	}

	AddStringFlag(cmd, "access-key-id", "i", "", "阿里云 Access Key ID", false)
	AddStringFlag(cmd, "access-key-secret", "k", "", "阿里云 Access Key Secret", false)
	AddStringFlag(cmd, "token", "t", "", "阿里云 Token", false)
	AddStringFlag(cmd, "role-arn", "r", "", "阿里云角色 ARN", false)
	AddStringFlag(cmd, "login-url", "l", "https://signin.aliyun.com/federation", "阿里云联合身份验证登录URL", false)
	AddStringFlag(cmd, "destination", "d", "https://console.aliyun.com", "登录成功后跳转的 URL", false)
	AddBoolFlag(cmd, "auto-login", "a", true, "自动打开 URL", false)

	// 关闭参数自动排序
	cmd.Flags().SortFlags = false

	return cmd
}

func runAlibaba() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// 获取命令行参数
		cfg := config.New()
		provider := cfg.GetProvider().GetAlibaba()

		accessKeyId, _ := cmd.Flags().GetString("access-key-id")
		accessKeySecret, _ := cmd.Flags().GetString("access-key-secret")
		token, _ := cmd.Flags().GetString("token")
		roleArn, _ := cmd.Flags().GetString("role-arn")
		loginUrl, _ := cmd.Flags().GetString("login-url")
		destination, _ := cmd.Flags().GetString("destination")

		// 如果命令行参数为空，尝试从配置文件读取
		if accessKeyId == "" || accessKeySecret == "" {
			// 从配置文件获取凭证信息
			if credential := provider.GetCredential(); credential != nil {
				if accessKeyId == "" {
					accessKeyId = credential.SecretId
				}
				if accessKeySecret == "" {
					accessKeySecret = credential.SecretKey
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
		if loginUrl == "" {
			loginUrl = provider.GetLoginUrl()
		}

		// 创建 Console 实例
		opts := &console.Options{
			Mode:     "cli",
			Provider: console.ProviderAlibaba,
		}

		c, err := console.New(opts)
		if err != nil {
			fmt.Printf("创建 Console 实例失败: %v\n", err)
			os.Exit(1)
		}

		// 创建登录选项
		loginOpts := console.NewLoginOptions(accessKeyId, accessKeySecret, token, roleArn, destination, loginUrl)

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
