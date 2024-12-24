package service

import (
	"fmt"

	"github.com/sg4i/cloud-console/internal/logic"
	"github.com/sg4i/cloud-console/internal/provider"
)

// ConsoleLoginOptions 定义控制台登录所需的基础参数
type ConsoleLoginOptions interface {
	GetCredential() *provider.Credential
	GetLoginURL() string
	GetDestination() string
	GetRoleArn() string
	GetAssumeRoleOptions() *provider.AssumeRoleOptions
}

// BaseLoginOptions 提供基础的登录选项实现
type BaseLoginOptions struct {
	Credential        *provider.Credential
	LoginURL          string
	Destination       string
	RoleArn           string
	AssumeRoleOptions *provider.AssumeRoleOptions
}

func (b *BaseLoginOptions) GetCredential() *provider.Credential { return b.Credential }
func (b *BaseLoginOptions) GetLoginURL() string                 { return b.LoginURL }
func (b *BaseLoginOptions) GetDestination() string              { return b.Destination }
func (b *BaseLoginOptions) GetRoleArn() string                  { return b.RoleArn }
func (b *BaseLoginOptions) GetAssumeRoleOptions() *provider.AssumeRoleOptions {
	return b.AssumeRoleOptions
}

// Provider 定义云服务提供商类型
type Provider string

const (
	ProviderAWS     Provider = "aws"
	ProviderAlibaba Provider = "alibaba"
	ProviderTencent Provider = "tencent"
)

// GenerateConsoleLoginURL 根据不同云厂商生成控制台登录URL
func GenerateConsoleLoginURL(provider Provider, opts ConsoleLoginOptions) (string, error) {
	switch provider {
	case ProviderTencent:
		tencentOpts := &logic.TencentLoginOptions{
			Credential:        opts.GetCredential(),
			SUrl:              opts.GetDestination(),
			RoleArn:           opts.GetRoleArn(),
			AssumeRoleOptions: opts.GetAssumeRoleOptions(),
		}
		return logic.GenerateTencentRoleLoginURL(tencentOpts)

	case ProviderAlibaba:
		alibabaOpts := &logic.AlibabaLoginOptions{
			Credential:        opts.GetCredential(),
			LoginURL:          opts.GetLoginURL(),
			Destination:       opts.GetDestination(),
			RoleArn:           opts.GetRoleArn(),
			AssumeRoleOptions: opts.GetAssumeRoleOptions(),
		}
		return logic.GenerateAlibabaRoleLoginURL(alibabaOpts)

	case ProviderAWS:
		awsOpts := &logic.AwsLoginOptions{
			Credential:        opts.GetCredential(),
			LoginURL:          opts.GetLoginURL(),
			Destination:       opts.GetDestination(),
			RoleArn:           opts.GetRoleArn(),
			AssumeRoleOptions: opts.GetAssumeRoleOptions(),
		}
		return logic.GenerateAwsRoleLoginURL(awsOpts)

	default:
		return "", fmt.Errorf("不支持的云厂商: %s", provider)
	}
}
