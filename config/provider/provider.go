package provider

import (
	internalProvider "github.com/sg4i/cloud-console/internal/provider"
)

type Provider interface {
	GetCredential() *internalProvider.Credential
	GetRoleArn() string
	GetLoginUrl() string
	GetDestination() string
}

type ProviderConfig struct {
	Tencent Provider
	Alibaba Provider
	Aws     Provider
}

func New() *ProviderConfig {
	return &ProviderConfig{
		Tencent: NewTencent(),
		Alibaba: NewAlibaba(),
		Aws:     NewAws(),
	}
}

// GetTencent 获取腾讯云配置
func (p *ProviderConfig) GetTencent() Provider {
	return p.Tencent
}

// GetAlibaba 获取阿里云配置
func (p *ProviderConfig) GetAlibaba() Provider {
	return p.Alibaba
}

// GetAws 获取AWS配置
func (p *ProviderConfig) GetAws() Provider {
	return p.Aws
}
