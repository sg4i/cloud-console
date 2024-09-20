package tencent

import (
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

// NewCommonClient 创建一个通用的腾讯云客户端
func NewCommonClient(secretId, secretKey, endpoint, region string) *common.Client {
	// 创建凭证
	credential := common.NewCredential(secretId, secretKey)
	logger.Log.Debug("已创建腾讯云凭证")

	// 创建客户端配置
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = endpoint
	logger.Log.WithField("endpoint", endpoint).Debug("已设置客户端配置")

	// 创建并返回 CommonClient
	client := common.NewCommonClient(credential, region, cpf)
	logger.Log.WithField("region", region).Debug("已创建腾讯云通用客户端")

	return client
}
