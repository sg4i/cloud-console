package cli

import (
	"fmt"
	"os"

	"github.com/sg4i/cloud-console/config"
	"github.com/sg4i/cloud-console/pkg/console"
	"github.com/spf13/cobra"
)

func NewAWSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws",
		Short: "生成AWS角色登录 URL",
		Run:   runAWS(),
	}

	AddStringFlag(cmd, "access-key-id", "i", "", "AWS Access Key ID", false)
	AddStringFlag(cmd, "secret-access-key", "k", "", "AWS Secret Access Key", false)
	AddStringFlag(cmd, "token", "t", "", "AWS Session Token", false)
	AddStringFlag(cmd, "role-arn", "r", "", "AWS角色 ARN", false)
	AddStringFlag(cmd, "login-url", "l", "https://signin.aws.amazon.com/federation", "AWS联合身份验证登录URL", false)
	AddStringFlag(cmd, "destination", "d", "https://console.aws.amazon.com", "登录成功后跳转的 URL", false)
	AddBoolFlag(cmd, "auto-login", "a", true, "自动打开 URL", false)

	// 关闭参数自动排序
	cmd.Flags().SortFlags = false

	return cmd
}

func runAWS() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// 获取命令行参数
		cfg := config.New()
		provider := cfg.GetProvider().GetAws()

		accessKeyId, _ := cmd.Flags().GetString("access-key-id")
		secretAccessKey, _ := cmd.Flags().GetString("secret-access-key")
		token, _ := cmd.Flags().GetString("token")
		roleArn, _ := cmd.Flags().GetString("role-arn")
		loginUrl, _ := cmd.Flags().GetString("login-url")
		destination, _ := cmd.Flags().GetString("destination")

		// 如果命令行参数为空，尝试从配置文件读取
		if accessKeyId == "" || secretAccessKey == "" {
			// 从配置文件获取凭证信息
			if credential := provider.GetCredential(); credential != nil {
				if accessKeyId == "" {
					accessKeyId = credential.SecretId
				}
				if secretAccessKey == "" {
					secretAccessKey = credential.SecretKey
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
			Provider: console.ProviderAws,
		}

		c, err := console.New(opts)
		if err != nil {
			fmt.Printf("创建 Console 实例失败: %v\n", err)
			os.Exit(1)
		}

		// 创建登录选项
		loginOpts := console.NewLoginOptions(accessKeyId, secretAccessKey, token, roleArn, destination, loginUrl)

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
