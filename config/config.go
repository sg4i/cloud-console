package config

import (
	"github.com/sg4i/cloud-console/config/grpc"
	"github.com/sg4i/cloud-console/config/provider"
	"github.com/sg4i/cloud-console/internal/utils"
)

// Config 总配置结构体
type Config struct {
	Provider *provider.ProviderConfig
	RPC      *grpc.RPCConfig
}

// New 创建新的配置实例
func New() *Config {
	// 加载配置文件
	utils.LoadConfig()

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
func (c *Config) GetRPC() *grpc.RPCConfig {
	return c.RPC
}
