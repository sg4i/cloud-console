package provider

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sirupsen/logrus"
)

// NewAlibabaStsClient 创建一个阿里云 STS 客户端
func NewAlibabaStsClient(accessKeyId, accessKeySecret string) (*sts20150401.Client, error) {
	openapiConfig := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("sts.cn-hangzhou.aliyuncs.com"),
	}

	client, err := sts20150401.NewClient(openapiConfig)
	if err != nil {
		logger.Log.WithError(err).Error("创建阿里云 STS 客户端失败")
		return nil, fmt.Errorf("创建阿里云 STS 客户端失败: %w", err)
	}

	logger.Log.Debug("已创建阿里云 STS 客户端")
	return client, nil
}

func AlibabaAssumeRole(accessKeyId, accessKeySecret string, opts *AssumeRoleOptions) (*Credential, error) {
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

	// 创建 STS 客户端
	client, err := NewAlibabaStsClient(accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("AssumeRole 请求失败: %w", err)
	}

	// 返回临时凭证
	logger.Log.Info("成功获取临时凭证")
	logger.Log.WithFields(logrus.Fields{
		"AccessKeyId":     tea.StringValue(response.Body.Credentials.AccessKeyId),
		"AccessKeySecret": tea.StringValue(response.Body.Credentials.AccessKeySecret),
		"SecurityToken":   tea.StringValue(response.Body.Credentials.SecurityToken),
		"Expiration":      tea.StringValue(response.Body.Credentials.Expiration),
	}).Debug("临时凭证详情")

	return &Credential{
		SecretId:  tea.StringValue(response.Body.Credentials.AccessKeyId),
		SecretKey: tea.StringValue(response.Body.Credentials.AccessKeySecret),
		Token:     tea.StringValue(response.Body.Credentials.SecurityToken),
	}, nil
}
