package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/spf13/viper"
)

// Profile 结构体包含所有需要的字段
type Profile struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	RoleArn         string
	Destination     string
	LoginUrl        string
}

func LoadProfile(cmdAccessKeyId, cmdAccessKeySecret, cmdSecurityToken, cmdRoleArn, cmdDestination, cmdLoginUrl string) (*Profile, error) {
	v := viper.New()

	// 设置配置键
	v.SetEnvPrefix("ALIBABA_CLOUD")
	v.BindEnv("ACCESS_KEY_ID")
	v.BindEnv("ACCESS_KEY_SECRET")
	v.BindEnv("SECURITY_TOKEN")
	v.BindEnv("ROLE_ARN")
	v.BindEnv("DESTINATION")
	v.BindEnv("LOGIN_URL")

	// 设置配置文件
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Log.WithError(err).Error("无法获取用户主目录")
		return nil, fmt.Errorf("无法获取用户主目录: %w", err)
	}

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(homeDir, ".cloud-console", "alibaba"))

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Log.WithError(err).Error("读取配置文件失败")
		}
	}

	// 设置默认值
	v.SetDefault("accessKeyId", "")
	v.SetDefault("accessKeySecret", "")
	v.SetDefault("securityToken", "")
	v.SetDefault("roleArn", "")
	v.SetDefault("destination", "")
	v.SetDefault("loginUrl", "https://signin.aliyun.com/federation")

	// 命令行参数（最高优先级）
	if cmdAccessKeyId != "" {
		v.Set("accessKeyId", cmdAccessKeyId)
	}
	if cmdAccessKeySecret != "" {
		v.Set("accessKeySecret", cmdAccessKeySecret)
	}
	if cmdSecurityToken != "" {
		v.Set("securityToken", cmdSecurityToken)
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
		AccessKeySecret: v.GetString("accessKeySecret"),
		SecurityToken:   v.GetString("securityToken"),
		RoleArn:         v.GetString("roleArn"),
		Destination:     v.GetString("destination"),
		LoginUrl:        v.GetString("loginUrl"),
	}

	// 验证必要的凭证信息
	if profile.AccessKeyId == "" || profile.AccessKeySecret == "" {
		return nil, fmt.Errorf("缺少必要的凭证信息")
	}

	// 新增验证：当 SecurityToken 为空时，RoleArn 不能为空
	if profile.SecurityToken == "" && profile.RoleArn == "" {
		return nil, fmt.Errorf("使用永久密钥时需指定角色 ARN")
	}

	return profile, nil
}
