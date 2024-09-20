package engine

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sg4i/cloud-console/internal/config"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sg4i/cloud-console/internal/tencent"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateTencentRoleLoginURL(cmdSecretId, cmdSecretKey, cmdToken, cmdSUrl string) (string, error) {
	cred, err := config.LoadCredentials(cmdSecretId, cmdSecretKey, cmdToken)
	if err != nil {
		return "", fmt.Errorf("加载凭证失败: %w", err)
	}

	algorithm := "sha1"
	
	sUrl := cmdSUrl
	if sUrl == "" {
		sUrl = "https://console.cloud.tencent.com"
	}

	nonce := rand.Intn(900000) + 100000
	timestamp := time.Now().Unix()

	url, err := tencent.GenerateRoleLoginURL(cred.SecretId, cred.SecretKey, cred.Token, nonce, timestamp, algorithm, sUrl)
	if err != nil {
		logger.Log.WithError(err).Error("生成腾讯云角色登录 URL 失败")
		return "", fmt.Errorf("生成腾讯云角色登录 URL 失败: %w", err)
	}

	logger.Log.WithField("url", url).Info("成功生成腾讯云角色登录 URL")
	return url, nil
}
