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

	// 创建客户端配置
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = endpoint

	// 创建并返回 CommonClient
	client := common.NewCommonClient(credential, region, cpf)

	return client
}

func TencentAssumeRole(secretId, secretKey, roleArn string) (*Credential, error) {

	opts := &AssumeRoleOptions{
		RoleSessionName: DefaultRoleSessionName,
		DurationSeconds: DefaultDurationSeconds,
		RoleArn:         roleArn,
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
	logger.Log.Infof("调用AssumeRole成功获取角色 %s 临时凭证", opts.RoleArn)
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
