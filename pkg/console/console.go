package console

import (
	"fmt"

	"github.com/sg4i/cloud-console/internal/provider"
	"github.com/sg4i/cloud-console/internal/service"
)

// Provider 定义云服务提供商类型
type Provider string

const (
	ProviderTencent Provider = "tencent"
	ProviderAlibaba Provider = "ali"
	ProviderAws     Provider = "aws"
)

// Options 定义Console的配置选项
type Options struct {
	// 服务模式：cli, grpc
	Mode string
	// 云服务提供商
	Provider Provider
}

// LoginOptions 定义登录选项
type LoginOptions struct {
	service.BaseLoginOptions
}

// Console 定义控制台核心结构
type Console struct {
	opts *Options
}

// New 创建新的Console实例
func New(opts *Options) (*Console, error) {
	if opts == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}

	return &Console{
		opts: opts,
	}, nil
}

// GetLoginURL 获取云服务控制台登录URL
func (c *Console) GetLoginURL(opts *LoginOptions) (string, error) {
	return service.GenerateConsoleLoginURL(service.Provider(c.opts.Provider), opts)
}

// NewLoginOptions 创建登录选项
func NewLoginOptions(secretId, secretKey, token, roleArn, destination, loginURL string) *LoginOptions {
	return &LoginOptions{
		BaseLoginOptions: service.BaseLoginOptions{
			Credential: &provider.Credential{
				SecretId:  secretId,
				SecretKey: secretKey,
				Token:     token,
			},
			RoleArn:     roleArn,
			Destination: destination,
			LoginURL:    loginURL,
		},
	}
}
