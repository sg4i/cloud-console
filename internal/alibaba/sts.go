package alibaba

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sg4i/cloud-console/internal/alibaba/config"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sirupsen/logrus"
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

func AssumeRole(accessKeyId, accessKeySecret string, opts *AssumeRoleOptions) (config.Credential, error) {
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

	// 创建 STS 客户端
	openapiConfig := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("sts.cn-hangzhou.aliyuncs.com"),
	}
	client, err := sts20150401.NewClient(openapiConfig)
	if err != nil {
		logger.Log.WithError(err).Error("创建 STS 客户端失败")
		return config.Credential{}, fmt.Errorf("创建 STS 客户端失败: %w", err)
	}

	// 创建 AssumeRole 请求
	request := &sts20150401.AssumeRoleRequest{
		RoleArn:         tea.String(opts.RoleArn),
		RoleSessionName: tea.String(opts.RoleSessionName),
		DurationSeconds: tea.Int64(int64(opts.DurationSeconds)),
	}

	// 发送请求
	response, err := client.AssumeRole(request)
	if err != nil {
		logger.Log.WithError(err).Error("AssumeRole 请求失败")
		return config.Credential{}, fmt.Errorf("AssumeRole 请求失败: %w", err)
	}

	// 返回临时凭证
	logger.Log.Info("成功获取临时凭证")
	logger.Log.WithFields(logrus.Fields{
		"AccessKeyId":     tea.StringValue(response.Body.Credentials.AccessKeyId),
		"AccessKeySecret": tea.StringValue(response.Body.Credentials.AccessKeySecret),
		"SecurityToken":   tea.StringValue(response.Body.Credentials.SecurityToken),
		"Expiration":      tea.StringValue(response.Body.Credentials.Expiration),
	}).Debug("临时凭证详情")
	return config.Credential{
		AccessKeyId:     tea.StringValue(response.Body.Credentials.AccessKeyId),
		AccessKeySecret: tea.StringValue(response.Body.Credentials.AccessKeySecret),
		SecurityToken:   tea.StringValue(response.Body.Credentials.SecurityToken),
	}, nil
}
