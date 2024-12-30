package config

import (
	"github.com/sg4i/cloud-console/config/grpc"
	"github.com/sg4i/cloud-console/config/provider"
	"github.com/sg4i/cloud-console/internal/utils"
)

// Config 总配置结构体
type Config struct {
	Provider *provider.ProviderConfig
	RPC      *grpc.GRPCConfig
}

// New 创建新的配置实例，configFile 为配置文件的完整路径（包含文件名）
func New(configFile string) *Config {
	// 加载配置文件
	utils.LoadConfig(configFile)

	return &Config{
		Provider: provider.New(),
		RPC:      grpc.New(),
	}
}

// GetProvider 获取云服务提供商配置
func (c *Config) GetProvider() *provider.ProviderConfig {
	return c.Provider
}

// GetRPC 获取RPC配置
func (c *Config) GetRPC() *grpc.GRPCConfig {
	return c.RPC
}
