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

	return client, nil
}

func AlibabaAssumeRole(accessKeyId, accessKeySecret, roleArn string) (*Credential, error) {
	opts := &AssumeRoleOptions{
		RoleSessionName: DefaultRoleSessionName,
		DurationSeconds: DefaultDurationSeconds,
		RoleArn: roleArn,
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
	logger.Log.Infof("调用AssumeRole成功获取角色%s的临时凭证", opts.RoleArn)
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
