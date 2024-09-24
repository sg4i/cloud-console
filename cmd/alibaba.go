package cmd

import (
	"fmt"

	"github.com/sg4i/cloud-console/internal/engine"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sg4i/cloud-console/internal/utils"
	"github.com/spf13/cobra"
)

var (
	alibabaAccessKeyId     string
	alibabaAccessKeySecret string
	alibabaSecurityToken   string
	alibabaLoginURL        string
	alibabaDestination     string
	alibabaRoleArn         string
)

var alibabaCmd = &cobra.Command{
	Use:   "alibaba",
	Short: "生成阿里云角色登录 URL",
	Long:  `生成阿里云角色登录 URL，用于在浏览器中直接访问阿里云控制台。`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := engine.GenerateAlibabaRoleLoginURL(
			alibabaAccessKeyId,
			alibabaAccessKeySecret,
			alibabaSecurityToken,
			alibabaRoleArn,
			alibabaDestination,
			alibabaLoginURL,
		)
		if err != nil {
			logger.Log.WithError(err).Error("生成阿里云角色登录 URL 失败")
			fmt.Println("生成阿里云角色登录 URL 失败:", err)
			return
		}
		fmt.Println("阿里云角色登录 URL:", url)

		if autoLogin {
			err = utils.OpenURL(url)
			if err != nil {
				logger.Log.WithError(err).Error("自动打开 URL 失败")
				fmt.Println("自动打开 URL 失败:", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(alibabaCmd)

	alibabaCmd.Flags().StringVar(&alibabaAccessKeyId, "access-key-id", "", "阿里云 AccessKeyId")
	alibabaCmd.Flags().StringVar(&alibabaAccessKeySecret, "access-key-secret", "", "阿里云 AccessKeySecret")
	alibabaCmd.Flags().StringVar(&alibabaSecurityToken, "security-token", "", "阿里云 SecurityToken（可选）")
	alibabaCmd.Flags().StringVar(&alibabaLoginURL, "login-url", "https://signin.aliyun.com/federation", "登录 URL（可选）")
	alibabaCmd.Flags().StringVar(&alibabaDestination, "destination", "https://home.console.aliyun.com", "目标 URL（可选）")
	alibabaCmd.Flags().StringVar(&alibabaRoleArn, "role-arn", "", "阿里云角色 ARN（可选）")

	// alibabaCmd.MarkFlagRequired("access-key-id")
	// alibabaCmd.MarkFlagRequired("access-key-secret")
}
