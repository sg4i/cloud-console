package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sg4i/cloud-console/internal/engine"
	"github.com/sg4i/cloud-console/internal/logger"
)

var (
	tencentSecretId  string
	tencentSecretKey string
	tencentToken     string
	tencentSUrl      string
	tencentRoleArn   string
)

var tencentCmd = &cobra.Command{
	Use:   "tencent",
	Short: "Generate Tencent Cloud role login URL",
	Run: func(cmd *cobra.Command, args []string) {
		url, err := engine.GenerateTencentRoleLoginURL(tencentSecretId, tencentSecretKey, tencentToken, tencentSUrl, tencentRoleArn)
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
	tencentCmd.Flags().StringVar(&tencentSecretId, "secret-id", "", "腾讯云 Secret ID")
	tencentCmd.Flags().StringVar(&tencentSecretKey, "secret-key", "", "腾讯云 Secret Key")
	tencentCmd.Flags().StringVar(&tencentToken, "token", "", "腾讯云 Token")
	tencentCmd.Flags().StringVar(&tencentSUrl, "s-url", "https://console.cloud.tencent.com", "登录成功后跳转的 URL")
	tencentCmd.Flags().StringVar(&tencentRoleArn, "role-arn", "", "腾讯云角色 ARN")
}
