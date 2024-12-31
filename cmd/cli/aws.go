package cli

import (
	"os"
	"path/filepath"

	"github.com/sg4i/cloud-console/cmd/common"
	"github.com/sg4i/cloud-console/config"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sg4i/cloud-console/internal/utils"
	"github.com/sg4i/cloud-console/pkg/console"
	"github.com/spf13/cobra"
)

func NewAWSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aws",
		Short: "生成AWS角色登录 URL",
		Run:   runAWS(),
	}

	// 设置默认配置文件路径
	defaultConfig := filepath.Join(".", "config.yml")

	common.AddStringFlag(cmd, "config", "c", defaultConfig, "配置文件路径", false)
	common.AddStringFlag(cmd, "access-key-id", "i", "", "AWS Access Key ID", false)
	common.AddStringFlag(cmd, "secret-access-key", "k", "", "AWS Secret Access Key", false)
	common.AddStringFlag(cmd, "token", "t", "", "AWS Session Token", false)
	common.AddStringFlag(cmd, "role-arn", "r", "", "AWS角色 ARN", false)
	common.AddStringFlag(cmd, "login-url", "l", "https://signin.aws.amazon.com/federation", "AWS联合身份验证登录URL", false)
	common.AddStringFlag(cmd, "destination", "d", "https://console.aws.amazon.com", "登录成功后跳转的 URL", false)
	common.AddBoolFlag(cmd, "auto-login", "a", true, "自动打开 URL", false)

	// 关闭参数自动排序
	cmd.Flags().SortFlags = false

	return cmd
}

func runAWS() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")
		accessKeyId, _ := cmd.Flags().GetString("access-key-id")
		secretAccessKey, _ := cmd.Flags().GetString("secret-access-key")
		token, _ := cmd.Flags().GetString("token")
		roleArn, _ := cmd.Flags().GetString("role-arn")
		loginUrl, _ := cmd.Flags().GetString("login-url")
		destination, _ := cmd.Flags().GetString("destination")
		autoLogin, _ := cmd.Flags().GetBool("auto-login")

		cfg := config.New(configFile)
		provider := cfg.GetProvider().GetAws()

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

		// 校验必要参数
		if accessKeyId == "" || secretAccessKey == "" {
			logger.Log.Error("Access Key ID 和 Secret Access Key 不能为空")
			os.Exit(1)
		}

		// 创建 Console 实例
		opts := &console.Options{
			Mode:     "cli",
			Provider: console.ProviderAws,
		}

		c, err := console.New(opts)
		if err != nil {
			logger.Log.WithError(err).Error("创建 Console 实例失败")
			os.Exit(1)
		}

		// 创建登录选项
		loginOpts := console.NewLoginOptions(accessKeyId, secretAccessKey, token, roleArn, destination, loginUrl)

		// 获取登录 URL
		url, err := c.GetLoginURL(loginOpts)
		if err != nil {
			logger.Log.WithError(err).Error("获取登录 URL 失败")
			os.Exit(1)
		}

		logger.Log.Info("AWS 控制台登录链接: " + url)
		if autoLogin {
			err = utils.OpenURL(url)
			if err != nil {
				logger.Log.WithError(err).Error("自动打开 URL 失败")
			}
		}
	}
}
