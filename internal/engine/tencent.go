package engine

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sg4i/cloud-console/internal/config"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sg4i/cloud-console/internal/tencent"
)

func GenerateTencentRoleLoginURL(cmdSecretId, cmdSecretKey, cmdToken, cmdSUrl, roleArn string) (string, error) {
	cred, err := config.LoadCredentials(cmdSecretId, cmdSecretKey, cmdToken)
	if err != nil {
		return "", fmt.Errorf("加载凭证失败: %w", err)
	}
	// 如果没有有效的 token，则使用 AssumeRole 获取临时密钥
	if cred.Token == "" {
		if roleArn == "" {
			return "", fmt.Errorf("Token为空情况下，角色 ARN 不能为空")
		}
		tempCred, err := tencent.AssumeRole(cred.SecretId, cred.SecretKey, &tencent.AssumeRoleOptions{RoleArn: roleArn})
		if err != nil {
			return "", fmt.Errorf("获取临时密钥失败: %w", err)
		}
		cred = &tempCred
	}

	sUrl := cmdSUrl
	if sUrl == "" {
		sUrl = "https://console.cloud.tencent.com"
	}

	algorithm := "sha1"
	nonce := rand.Intn(900000) + 100000
	timestamp := time.Now().Unix()

	url, err := tencent.GenerateRoleLoginURL(cred.SecretId, cred.SecretKey, cred.Token, nonce, timestamp, algorithm, sUrl)
	if err != nil {
		logger.Log.WithError(err).Error("生成腾讯云角色登录 URL 失败")
		return "", fmt.Errorf("生成腾讯云角色登录 URL 失败: %w", err)
	}

	logger.Log.WithField("url", url).Debug("成功生成腾讯云角色登录 URL")
	return url, nil
}
