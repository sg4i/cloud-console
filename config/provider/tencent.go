package provider

import (
	internalProvider "github.com/sg4i/cloud-console/internal/provider"
	"github.com/sg4i/cloud-console/internal/utils"
)

type TencentConfig struct {
	Credential  *internalProvider.Credential
	RoleArn     string
	LoginUrl    string
	Destination string
}

func NewTencent() *TencentConfig {
	utils.LoadConfig()

	config := &TencentConfig{
		Credential: &internalProvider.Credential{},
	}

	// 读取认证信息
	config.Credential.SecretId = utils.GetString("provider.tencent.credential.secretId")
	config.Credential.SecretKey = utils.GetString("provider.tencent.credential.secretKey")
	config.Credential.Token = utils.GetString("provider.tencent.credential.token")

	// 读取其他配置
	config.RoleArn = utils.GetString("provider.tencent.roleArn")
	config.LoginUrl = utils.GetString("provider.tencent.loginUrl")
	config.Destination = utils.GetString("provider.tencent.destination")

	return config
}

// GetCredential 获取认证信息
func (c *TencentConfig) GetCredential() *internalProvider.Credential {
	return c.Credential
}

// GetRoleArn 获取角色ARN
func (c *TencentConfig) GetRoleArn() string {
	return c.RoleArn
}

// GetLoginUrl 获取登录URL
func (c *TencentConfig) GetLoginUrl() string {
	return c.LoginUrl
}

// GetDestination 获取目标地址
func (c *TencentConfig) GetDestination() string {
	return c.Destination
}
