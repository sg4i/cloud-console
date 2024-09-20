package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sg4i/cloud-console/internal/engine"
	"github.com/sg4i/cloud-console/internal/logger"
)

var (
	secretId  string
	secretKey string
	token     string
	sUrl      string
	roleArn   string
)

var tencentCmd = &cobra.Command{
	Use:   "tencent",
	Short: "Generate Tencent Cloud role login URL",
	Run: func(cmd *cobra.Command, args []string) {
		url, err := engine.GenerateTencentRoleLoginURL(secretId, secretKey, token, sUrl, roleArn)
		if err != nil {
			logger.Log.WithError(err).Error("生成腾讯云角色登录 URL 失败")
			return
		}
		fmt.Println("腾讯云角色登录 URL:")
		fmt.Println(url)
	},
}

func init() {
	rootCmd.AddCommand(tencentCmd)
	tencentCmd.Flags().StringVar(&secretId, "secret-id", "", "腾讯云 Secret ID")
	tencentCmd.Flags().StringVar(&secretKey, "secret-key", "", "腾讯云 Secret Key")
	tencentCmd.Flags().StringVar(&token, "token", "", "腾讯云 Token")
	tencentCmd.Flags().StringVar(&sUrl, "s-url", "https://console.cloud.tencent.com", "登录成功后跳转的 URL")
	tencentCmd.Flags().StringVar(&roleArn, "role-arn", "", "腾讯云角色 ARN")
}
