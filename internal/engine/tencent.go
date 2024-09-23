package engine

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sg4i/cloud-console/internal/tencent/config"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sg4i/cloud-console/internal/tencent"
)

func GenerateTencentRoleLoginURL(cmdSecretId, cmdSecretKey, cmdToken, cmdSUrl, cmdRoleArn string) (string, error) {
	profile, err := config.LoadProfile(cmdSecretId, cmdSecretKey, cmdToken, cmdRoleArn, cmdSUrl)
	if err != nil {
		return "", fmt.Errorf("加载配置文件失败: %w", err)
	}

	// 如果没有有效的 token，则使用 AssumeRole 获取临时密钥
	if profile.Token == "" {
		tempCred, err := tencent.AssumeRole(profile.SecretId, profile.SecretKey, &tencent.AssumeRoleOptions{RoleArn: profile.RoleArn})
		if err != nil {
			return "", fmt.Errorf("获取临时密钥失败: %w", err)
		}
		profile.SecretId = tempCred.SecretId
		profile.SecretKey = tempCred.SecretKey
		profile.Token = tempCred.Token
	}

	algorithm := "sha1"
	nonce := rand.Intn(900000) + 100000
	timestamp := time.Now().Unix()

	url, err := tencent.GenerateRoleLoginURL(profile.SecretId, profile.SecretKey, profile.Token, nonce, timestamp, algorithm, profile.SURL)
	if err != nil {
		logger.Log.WithError(err).Error("生成腾讯云角色登录 URL 失败")
		return "", fmt.Errorf("生成腾讯云角色登录 URL 失败: %w", err)
	}

	logger.Log.WithField("url", url).Debug("成功生成腾讯云角色登录 URL")
	return url, nil
}
