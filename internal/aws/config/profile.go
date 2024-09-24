package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Profile 结构体包含所有需要的字段
type Profile struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
	RoleArn         string
	Destination     string
	LoginUrl        string
}

func LoadProfile(cmdAccessKeyId, cmdSecretAccessKey, cmdSessionToken, cmdRoleArn, cmdDestination, cmdLoginUrl string) (*Profile, error) {
	v := viper.New()

	// 设置配置键
	v.SetEnvPrefix("AWS")
	v.BindEnv("ACCESS_KEY_ID")
	v.BindEnv("SECRET_ACCESS_KEY")
	v.BindEnv("SESSION_TOKEN")
	v.BindEnv("ROLE_ARN")
	v.BindEnv("DESTINATION")
	v.BindEnv("LOGIN_URL")

	// 设置配置文件
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Log.WithError(err).Error("无法获取用户主目录")
		return nil, fmt.Errorf("无法获取用户主目录: %w", err)
	}

	v.SetConfigName("default")
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(homeDir, ".cloud-console", "aws"))

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Log.WithError(err).Error("读取配置文件失败")
		}
	}

	// 设置默认值
	v.SetDefault("accessKeyId", "")
	v.SetDefault("secretAccessKey", "")
	v.SetDefault("sessionToken", "")
	v.SetDefault("roleArn", "")
	v.SetDefault("destination", "https://console.aws.amazon.com")
	v.SetDefault("loginUrl", "https://signin.aws.amazon.com/federation")

	// 命令行参数（最高优先级）
	if cmdAccessKeyId != "" {
		v.Set("accessKeyId", cmdAccessKeyId)
	}
	if cmdSecretAccessKey != "" {
		v.Set("secretAccessKey", cmdSecretAccessKey)
	}
	if cmdSessionToken != "" {
		v.Set("sessionToken", cmdSessionToken)
	}
	if cmdRoleArn != "" {
		v.Set("roleArn", cmdRoleArn)
	}
	if cmdDestination != "" {
		v.Set("destination", cmdDestination)
	}
	if cmdLoginUrl != "" {
		v.Set("loginUrl", cmdLoginUrl)
	}

	// 构建 Profile
	profile := &Profile{
		AccessKeyId:     v.GetString("accessKeyId"),
		SecretAccessKey: v.GetString("secretAccessKey"),
		SessionToken:    v.GetString("sessionToken"),
		RoleArn:         v.GetString("roleArn"),
		Destination:     v.GetString("destination"),
		LoginUrl:        v.GetString("loginUrl"),
	}
	logger.Log.WithFields(logrus.Fields{
		"AccessKeyId":  profile.AccessKeyId,
		"SessionToken": profile.SessionToken,
		"RoleArn":      profile.RoleArn,
		"Destination":  profile.Destination,
		"LoginUrl":     profile.LoginUrl,
	}).Debug("Profile 配置信息")

	// 验证必要的凭证信息
	if profile.AccessKeyId == "" || profile.SecretAccessKey == "" {
		return nil, fmt.Errorf("缺少必要的凭证信息")
	}

	// 新增验证：当 SessionToken 为空时，RoleArn 不能为空
	if profile.SessionToken == "" && profile.RoleArn == "" {
		return nil, fmt.Errorf("使用永久密钥时需指定角色 ARN")
	}

	return profile, nil
}
