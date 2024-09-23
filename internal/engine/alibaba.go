package engine

import (
	"fmt"

	"github.com/sg4i/cloud-console/internal/alibaba/config"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sg4i/cloud-console/internal/alibaba"
)

func GenerateAlibabaRoleLoginURL(cmdAccessKeyId, cmdAccessKeySecret, cmdSecurityToken, cmdRoleArn, cmdDestination, cmdLoginUrl string) (string, error) {
	profile, err := config.LoadProfile(cmdAccessKeyId, cmdAccessKeySecret, cmdSecurityToken, cmdRoleArn, cmdDestination, cmdLoginUrl)
	if err != nil {
		return "", fmt.Errorf("加载配置文件失败: %w", err)
	}

	// 如果没有有效的 SecurityToken，则使用 AssumeRole 获取临时密钥
	if profile.SecurityToken == "" {
		tempCred, err := alibaba.AssumeRole(profile.AccessKeyId, profile.AccessKeySecret, profile.RoleArn)
		if err != nil {
			return "", fmt.Errorf("获取临时密钥失败: %w", err)
		}
		profile.AccessKeyId = tempCred.AccessKeyId
		profile.AccessKeySecret = tempCred.AccessKeySecret
		profile.SecurityToken = tempCred.SecurityToken
	}

	params := alibaba.RoleLoginParams{
		LoginURL:    profile.LoginUrl,
		Destination: profile.Destination,
		Credential: config.Credential{
			AccessKeyId:     profile.AccessKeyId,
			AccessKeySecret: profile.AccessKeySecret,
			SecurityToken:   profile.SecurityToken,
		},
	}

	url, err := alibaba.GenerateRoleLoginURL(params)
	if err != nil {
		logger.Log.WithError(err).Error("生成阿里云角色登录 URL 失败")
		return "", fmt.Errorf("生成阿里云角色登录 URL 失败: %w", err)
	}

	logger.Log.WithField("url", url).Debug("成功生成阿里云角色登录 URL")
	return url, nil
}
