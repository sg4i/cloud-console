package provider

import (
	"fmt"

	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
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

func TencentAssumeRole(secretId, secretKey string, opts *AssumeRoleOptions) (*Credential, error) {
	if opts == nil {
		opts = &AssumeRoleOptions{
			RoleSessionName: DefaultRoleSessionName,
			DurationSeconds: DefaultDurationSeconds,
		}
	}
	if opts.RoleArn == "" {
		return nil, fmt.Errorf("RoleArn 不能为空")
	}

	// 使用默认值填充未指定的选项
	if opts.RoleSessionName == "" {
		opts.RoleSessionName = DefaultRoleSessionName
	}
	if opts.DurationSeconds == 0 {
		opts.DurationSeconds = DefaultDurationSeconds
	}

	// 创建 CommonClient
	client := NewCommonClient(secretId, secretKey, "sts.tencentcloudapi.com", "ap-guangzhou")

	// 创建请求
	request := sts.NewAssumeRoleRequest()

	request.RoleArn = common.StringPtr(opts.RoleArn)
	request.RoleSessionName = common.StringPtr(opts.RoleSessionName)
	request.DurationSeconds = common.Uint64Ptr(opts.DurationSeconds)

	// 发送请求
	response := sts.NewAssumeRoleResponse()
	err := client.Send(request, response)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		logger.Log.WithError(err).Error("AssumeRole 请求失败")
		return nil, fmt.Errorf("AssumeRole 请求失败: %w", err)
	}
	if err != nil {
		logger.Log.WithError(err).Error("发送 AssumeRole 请求时发生错误")
		return nil, fmt.Errorf("发送 AssumeRole 请求时发生错误: %w", err)
	}

	// 返回临时凭证
	logger.Log.Info("成功获取临时凭证")
	logger.Log.WithFields(logrus.Fields{
		"TmpSecretId": *response.Response.Credentials.TmpSecretId,
		"Token":       *response.Response.Credentials.Token,
		"ExpiredTime": *response.Response.ExpiredTime,
	}).Debug("临时凭证详情")
	return &Credential{
		SecretId:  *response.Response.Credentials.TmpSecretId,
		SecretKey: *response.Response.Credentials.TmpSecretKey,
		Token:     *response.Response.Credentials.Token,
	}, nil
}
