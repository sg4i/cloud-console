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
	SecretId  string
	SecretKey string
	Token     string
	ARN       string
	SURL      string
}

func LoadProfile(cmdSecretId, cmdSecretKey, cmdToken, cmdARN, cmdSURL string) (*Profile, error) {
	v := viper.New()

	// 设置配置键
	v.SetEnvPrefix("TENCENTCLOUD")
	v.BindEnv("SECRET_ID")
	v.BindEnv("SECRET_KEY")
	v.BindEnv("TOKEN")
	v.BindEnv("ARN")
	v.BindEnv("SURL")

	// 设置配置文件
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Log.WithError(err).Error("无法获取用户主目录")
		return nil, fmt.Errorf("无法获取用户主目录: %w", err)
	}

	v.SetConfigName("default")
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Join(homeDir, ".cloud-console"))

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Log.WithError(err).Error("读取配置文件失败")
		}
	}

	// 设置默认值
	v.SetDefault("secretId", "")
	v.SetDefault("secretKey", "")
	v.SetDefault("token", "")
	v.SetDefault("arn", "")
	v.SetDefault("surl", "https://console.cloud.tencent.com")

	// 命令行参数（最高优先级）
	if cmdSecretId != "" {
		v.Set("secretId", cmdSecretId)
	}
	if cmdSecretKey != "" {
		v.Set("secretKey", cmdSecretKey)
	}
	if cmdToken != "" {
		v.Set("token", cmdToken)
	}
	if cmdARN != "" {
		v.Set("arn", cmdARN)
	}
	if cmdSURL != "" {
		v.Set("surl", cmdSURL)
	}

	// 构建 Profile
	profile := &Profile{
		SecretId:  v.GetString("secretId"),
		SecretKey: v.GetString("secretKey"),
		Token:     v.GetString("token"),
		ARN:       v.GetString("arn"),
		SURL:      v.GetString("surl"),
	}

	// 验证必要的凭证信息
	if profile.SecretId == "" || profile.SecretKey == "" {
		return nil, fmt.Errorf("缺少必要的凭证信息")
	}

	// 新增验证：当 token 为空时，arn 不能为空
	if profile.Token == "" && profile.ARN == "" {
		return nil, fmt.Errorf("使用永久密钥时需指定角色 ARN")
	}

	return profile, nil
}
