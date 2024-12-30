package provider

import (
	internalProvider "github.com/sg4i/cloud-console/internal/provider"
	"github.com/sg4i/cloud-console/internal/utils"
)

type AwsConfig struct {
	Credential  *internalProvider.Credential
	RoleArn     string
	LoginUrl    string
	Destination string
}

func NewAws() *AwsConfig {

	config := &AwsConfig{
		Credential: &internalProvider.Credential{},
	}

	// 读取认证信息
	config.Credential.SecretId = utils.GetString("provider.aws.credential.secretId")
	config.Credential.SecretKey = utils.GetString("provider.aws.credential.secretKey")
	config.Credential.Token = utils.GetString("provider.aws.credential.token")

	// 读取其他配置
	config.RoleArn = utils.GetString("provider.aws.roleArn")
	config.LoginUrl = utils.GetString("provider.aws.loginUrl")
	config.Destination = utils.GetString("provider.aws.destination")

	return config
}

// GetCredential 获取认证信息
func (c *AwsConfig) GetCredential() *internalProvider.Credential {
	return c.Credential
}

// GetRoleArn 获取角色ARN
func (c *AwsConfig) GetRoleArn() string {
	return c.RoleArn
}

// GetLoginUrl 获取登录URL
func (c *AwsConfig) GetLoginUrl() string {
	return c.LoginUrl
}

// GetDestination 获取目标地址
func (c *AwsConfig) GetDestination() string {
	return c.Destination
}
