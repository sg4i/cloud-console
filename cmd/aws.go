package cmd

import (
	"fmt"

	"github.com/sg4i/cloud-console/internal/engine"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sg4i/cloud-console/internal/utils"
	"github.com/spf13/cobra"
)

var (
	awsAccessKeyId     string
	awsSecretAccessKey string
	awsSessionToken    string
	awsRoleArn         string
	awsDestination     string
	awsLoginURL        string
)

var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "生成 AWS 角色登录 URL",
	Long:  `生成 AWS 角色登录 URL，用于在浏览器中直接访问 AWS 控制台。`,
	Run: func(cmd *cobra.Command, args []string) {
		url, err := engine.GenerateAWSRoleLoginURL(
			awsAccessKeyId,
			awsSecretAccessKey,
			awsSessionToken,
			awsRoleArn,
			awsDestination,
			awsLoginURL,
		)
		if err != nil {
			logger.Log.WithError(err).Error("生成 AWS 角色登录 URL 失败")
			fmt.Println("生成 AWS 角色登录 URL 失败:", err)
			return
		}
		fmt.Println("AWS 角色登录 URL:", url)

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
	rootCmd.AddCommand(awsCmd)

	awsCmd.Flags().StringVar(&awsAccessKeyId, "access-key-id", "", "AWS AccessKeyId")
	awsCmd.Flags().StringVar(&awsSecretAccessKey, "secret-access-key", "", "AWS SecretAccessKey")
	awsCmd.Flags().StringVar(&awsSessionToken, "session-token", "", "AWS SessionToken（可选）")
	awsCmd.Flags().StringVar(&awsLoginURL, "login-url", "https://signin.aws.amazon.com/federation", "登录 URL（可选）")
	awsCmd.Flags().StringVar(&awsDestination, "destination", "https://console.aws.amazon.com", "目标 URL（可选）")
	awsCmd.Flags().StringVar(&awsRoleArn, "role-arn", "", "AWS 角色 ARN（可选）")

	// awsCmd.MarkFlagRequired("access-key-id")
	// awsCmd.MarkFlagRequired("secret-access-key")
}
