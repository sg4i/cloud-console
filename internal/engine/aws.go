package engine

import (
	"fmt"

	"github.com/sg4i/cloud-console/internal/aws"
	"github.com/sg4i/cloud-console/internal/aws/config"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sirupsen/logrus"
)

func GenerateAWSRoleLoginURL(cmdAccessKeyId, cmdSecretAccessKey, cmdSessionToken, cmdRoleArn, cmdDestination, cmdLoginUrl string) (string, error) {
	profile, err := config.LoadProfile(cmdAccessKeyId, cmdSecretAccessKey, cmdSessionToken, cmdRoleArn, cmdDestination, cmdLoginUrl)
	if err != nil {
		return "", fmt.Errorf("加载配置文件失败: %w", err)
	}

	// 如果没有有效的 SessionToken，则使用 AssumeRole 获取临时密钥
	if profile.SessionToken == "" {
		opts := &aws.AssumeRoleOptions{
			RoleArn: profile.RoleArn,
		}
		tempCred, err := aws.AssumeRole(profile.AccessKeyId, profile.SecretAccessKey, profile.SessionToken, opts)
		if err != nil {
			return "", fmt.Errorf("获取临时密钥失败: %w", err)
		}
		profile.AccessKeyId = tempCred.AccessKeyId
		profile.SecretAccessKey = tempCred.SecretAccessKey
		profile.SessionToken = tempCred.SessionToken
	}

	params := aws.RoleLoginParams{
		LoginURL:    profile.LoginUrl,
		Destination: profile.Destination,
		Credential: config.Credential{
			AccessKeyId:     profile.AccessKeyId,
			SecretAccessKey: profile.SecretAccessKey,
			SessionToken:    profile.SessionToken,
		},
	}
	logger.Log.WithFields(logrus.Fields{
		"LoginURL":      params.LoginURL,
		"Destination":   params.Destination,
		"AccessKeyId":   params.Credential.AccessKeyId,
		"SessionToken":  params.Credential.SessionToken,
	}).Debug("生成 AWS 角色登录 URL 参数")

	url, err := aws.GenerateRoleLoginURL(params)
	if err != nil {
		logger.Log.WithError(err).Error("生成 AWS 角色登录 URL 失败")
		return "", fmt.Errorf("生成 AWS 角色登录 URL 失败: %w", err)
	}

	logger.Log.WithField("url", url).Debug("成功生成 AWS 角色登录 URL")
	return url, nil
}
