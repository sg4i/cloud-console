package provider

import (
	internalProvider "github.com/sg4i/cloud-console/internal/provider"
	"github.com/sg4i/cloud-console/internal/utils"
)

type AlibabaConfig struct {
	Credential  *internalProvider.Credential
	RoleArn     string
	LoginUrl    string
	Destination string
}

func NewAlibaba() *AlibabaConfig {
	utils.LoadConfig()

	config := &AlibabaConfig{
		Credential: &internalProvider.Credential{},
	}

	// 读取认证信息
	config.Credential.SecretId = utils.GetString("provider.alibaba.credential.secretId")
	config.Credential.SecretKey = utils.GetString("provider.alibaba.credential.secretKey")
	config.Credential.Token = utils.GetString("provider.alibaba.credential.token")

	// 读取其他配置
	config.RoleArn = utils.GetString("provider.alibaba.roleArn")
	config.LoginUrl = utils.GetString("provider.alibaba.loginUrl")
	config.Destination = utils.GetString("provider.alibaba.destination")

	return config
}

// GetCredential 获取认证信息
func (c *AlibabaConfig) GetCredential() *internalProvider.Credential {
	return c.Credential
}

// GetRoleArn 获取角色ARN
func (c *AlibabaConfig) GetRoleArn() string {
	return c.RoleArn
}

// GetLoginUrl 获取登录URL
func (c *AlibabaConfig) GetLoginUrl() string {
	return c.LoginUrl
}

// GetDestination 获取目标地址
func (c *AlibabaConfig) GetDestination() string {
	return c.Destination
}
