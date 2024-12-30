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

func NewTencentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tencent",
		Short: "生成腾讯云角色登录 URL",
		Run:   runTencent(),
	}

	// 设置默认配置文件路径
	defaultConfig := filepath.Join(".", "config.yml")

	common.AddStringFlag(cmd, "config", "c", defaultConfig, "配置文件路径", false)
	common.AddStringFlag(cmd, "secret-id", "i", "", "腾讯云 Secret ID", false)
	common.AddStringFlag(cmd, "secret-key", "k", "", "腾讯云 Secret Key", false)
	common.AddStringFlag(cmd, "token", "t", "", "腾讯云 Token", false)
	common.AddStringFlag(cmd, "role-arn", "r", "", "腾讯云角色 ARN", false)
	common.AddStringFlag(cmd, "destination", "d", "https://console.cloud.tencent.com", "登录成功后跳转的 URL", false)
	common.AddBoolFlag(cmd, "auto-login", "a", true, "自动打开 URL", false)

	// 关闭参数自动排序
	cmd.Flags().SortFlags = false

	return cmd
}

func runTencent() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		configFile, _ := cmd.Flags().GetString("config")
		secretId, _ := cmd.Flags().GetString("secret-id")
		secretKey, _ := cmd.Flags().GetString("secret-key")
		token, _ := cmd.Flags().GetString("token")
		roleArn, _ := cmd.Flags().GetString("role-arn")
		destination, _ := cmd.Flags().GetString("destination")
		autoLogin, _ := cmd.Flags().GetBool("auto-login")

		cfg := config.New(configFile)
		provider := cfg.GetProvider().GetTencent()

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
			logger.Log.Error("Secret ID 和 Secret Key 不能为空")
			os.Exit(1)
		}

		// 创建 Console 实例
		opts := &console.Options{
			Mode:     "cli",
			Provider: console.ProviderTencent,
		}

		c, err := console.New(opts)
		if err != nil {
			logger.Log.WithError(err).Error("创建 Console 实例失败")
			os.Exit(1)
		}

		// 创建登录选项
		loginOpts := console.NewLoginOptions(secretId, secretKey, token, roleArn, destination, "")

		// 获取登录 URL
		url, err := c.GetLoginURL(loginOpts)
		if err != nil {
			logger.Log.WithError(err).Error("获取登录 URL 失败")
			os.Exit(1)
		}

		logger.Log.Info(url)
		if autoLogin {
			err = utils.OpenURL(url)
			if err != nil {
				logger.Log.WithError(err).Error("自动打开 URL 失败")
			}
		}
	}
}
