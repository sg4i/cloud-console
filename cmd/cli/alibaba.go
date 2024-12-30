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

func NewAlibabaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alibaba",
		Short: "生成阿里云角色登录 URL",
		Run:   runAlibaba(),
	}

	// 设置默认配置文件路径
	defaultConfig := filepath.Join(".", "config.yml")

	common.AddStringFlag(cmd, "config", "c", defaultConfig, "配置文件路径", false)
	common.AddStringFlag(cmd, "access-key-id", "i", "", "阿里云 Access Key ID", false)
	common.AddStringFlag(cmd, "access-key-secret", "k", "", "阿里云 Access Key Secret", false)
	common.AddStringFlag(cmd, "token", "t", "", "阿里云 Token", false)
	common.AddStringFlag(cmd, "role-arn", "r", "", "阿里云角色 ARN", false)
	common.AddStringFlag(cmd, "login-url", "l", "https://signin.aliyun.com/federation", "阿里云联合身份验证登录URL", false)
	common.AddStringFlag(cmd, "destination", "d", "https://console.aliyun.com", "登录成功后跳转的 URL", false)
	common.AddBoolFlag(cmd, "auto-login", "a", true, "自动打开 URL", false)

	// 关闭参数自动排序
	cmd.Flags().SortFlags = false

	return cmd
}

func runAlibaba() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		configFile, _ := cmd.Flags().GetString("config")
		accessKeyId, _ := cmd.Flags().GetString("access-key-id")
		accessKeySecret, _ := cmd.Flags().GetString("access-key-secret")
		token, _ := cmd.Flags().GetString("token")
		roleArn, _ := cmd.Flags().GetString("role-arn")
		loginUrl, _ := cmd.Flags().GetString("login-url")
		destination, _ := cmd.Flags().GetString("destination")
		autoLogin, _ := cmd.Flags().GetBool("auto-login")

		cfg := config.New(configFile)
		provider := cfg.GetProvider().GetTencent()

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

		// 校验必要参数
		if accessKeyId == "" || accessKeySecret == "" {
			logger.Log.Error("Access Key ID 和 Access Key Secret 不能为空")
			os.Exit(1)
		}

		// 创建 Console 实例
		opts := &console.Options{
			Mode:     "cli",
			Provider: console.ProviderAlibaba,
		}

		c, err := console.New(opts)
		if err != nil {
			logger.Log.WithError(err).Error("创建 Console 实例失败")
			os.Exit(1)
		}

		// 创建登录选项
		loginOpts := console.NewLoginOptions(accessKeyId, accessKeySecret, token, roleArn, destination, loginUrl)

		// 获取登录 URL
		url, err := c.GetLoginURL(loginOpts)
		if err != nil {
			logger.Log.WithError(err).Error("获取登录 URL 失败")
			os.Exit(1)
		}

		if autoLogin {
			err = utils.OpenURL(url)
			if err != nil {
				logger.Log.WithError(err).Error("自动打开 URL 失败")
			}
		} else {
			logger.Log.Info(url)
		}
	}
}
