package tencent

import (
	"fmt"

	"github.com/sg4i/cloud-console/internal/tencent/config"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"
)

// AssumeRoleOptions 定义 AssumeRole 的可选参数
type AssumeRoleOptions struct {
	RoleArn         string
	RoleSessionName string
	DurationSeconds uint64
}

// DefaultAssumeRoleOptions 提供默认的 AssumeRole 选项
var DefaultAssumeRoleOptions = AssumeRoleOptions{
	RoleSessionName: "RoleSession",
	DurationSeconds: 3600,
}

func AssumeRole(secretId, secretKey string, opts *AssumeRoleOptions) (config.Credential, error) {
	if opts == nil {
		opts = &DefaultAssumeRoleOptions
	}
	if opts.RoleArn == "" {
		return config.Credential{}, fmt.Errorf("RoleArn 不能为空")
	}

	// 使用默认值填充未指定的选项
	if opts.RoleSessionName == "" {
		opts.RoleSessionName = DefaultAssumeRoleOptions.RoleSessionName
	}
	if opts.DurationSeconds == 0 {
		opts.DurationSeconds = DefaultAssumeRoleOptions.DurationSeconds
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
		return config.Credential{}, fmt.Errorf("AssumeRole 请求失败: %w", err)
	}
	if err != nil {
		logger.Log.WithError(err).Error("发送 AssumeRole 请求时发生错误")
		return config.Credential{}, fmt.Errorf("发送 AssumeRole 请求时发生错误: %w", err)
	}

	// 返回临时凭证
	logger.Log.Info("成功获取临时凭证")
	logger.Log.WithFields(logrus.Fields{
		"TmpSecretId": *response.Response.Credentials.TmpSecretId,
		"Token":       *response.Response.Credentials.Token,
		"ExpiredTime": *response.Response.ExpiredTime,
	}).Debug("临时凭证详情")
	return config.Credential{
		SecretId:  *response.Response.Credentials.TmpSecretId,
		SecretKey: *response.Response.Credentials.TmpSecretKey,
		Token:     *response.Response.Credentials.Token,
	}, nil
}
