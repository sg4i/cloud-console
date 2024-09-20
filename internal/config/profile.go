package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sg4i/cloud-console/internal/logger"
)

func LoadCredentials(cmdSecretId, cmdSecretKey, cmdToken string) (*Credential, error) {
	// 1. 命令行参数（最高优先级）
	if cmdSecretId != "" && cmdSecretKey != "" {
		return &Credential{
			SecretId:  cmdSecretId,
			SecretKey: cmdSecretKey,
			Token:     cmdToken,
		}, nil
	}

	// 2. 环境变量（次高优先级）
	envSecretId := os.Getenv("TENCENTCLOUD_SECRET_ID")
	envSecretKey := os.Getenv("TENCENTCLOUD_SECRET_KEY")
	envToken := os.Getenv("TENCENTCLOUD_TOKEN")
	if envSecretId != "" && envSecretKey != "" {
		return &Credential{
			SecretId:  envSecretId,
			SecretKey: envSecretKey,
			Token:     envToken,
		}, nil
	}

	// 3. 配置文件（最低优先级）
	return loadFromCredentialFile()
}

func loadFromCredentialFile() (*Credential, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Log.WithError(err).Error("无法获取用户主目录")
		return nil, fmt.Errorf("无法获取用户主目录: %w", err)
	}

	configPaths := []string{
		filepath.Join(homeDir, ".cloud-console", "default.credential"),
		filepath.Join(homeDir, ".tccli", "default.credential"),
	}

	for _, path := range configPaths {
		config, err := loadCredential(path)
		if err == nil {
			return config, nil
		}
	}

	return nil, fmt.Errorf("未找到有效的凭证信息")
}

func loadCredential(path string) (*Credential, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Log.WithError(err).WithField("path", path).Error("无法读取凭证文件")
		return nil, err
	}

	var config Credential
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.SecretId == "" || config.SecretKey == "" {
		return nil, fmt.Errorf("配置文件 %s 中缺少必要的凭证信息", path)
	}

	return &config, nil
}
