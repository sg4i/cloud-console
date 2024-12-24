package provider

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sirupsen/logrus"
)

// NewAwsStsClient 创建一个 AWS STS 客户端
func NewAwsStsClient(accessKeyId, secretAccessKey, sessionToken string) (*sts.STS, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyId, secretAccessKey, sessionToken),
		Region:      aws.String("us-east-1"),
	})
	if err != nil {
		logger.Log.WithError(err).Error("创建 AWS STS 客户端失败")
		return nil, fmt.Errorf("创建 AWS STS 客户端失败: %w", err)
	}

	logger.Log.Debug("已创建 AWS STS 客户端")
	return sts.New(sess), nil
}

func AwsAssumeRole(accessKeyId, secretAccessKey, sessionToken string, opts *AssumeRoleOptions) (*Credential, error) {
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
	client, err := NewAwsStsClient(accessKeyId, secretAccessKey, sessionToken)
	if err != nil {
		return nil, err
	}

	// 创建 AssumeRole 请求
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(opts.RoleArn),
		RoleSessionName: aws.String(opts.RoleSessionName),
		DurationSeconds: aws.Int64(int64(opts.DurationSeconds)),
	}

	// 发送请求
	result, err := client.AssumeRole(input)
	if err != nil {
		logger.Log.WithError(err).Error("AssumeRole 请求失败")
		return nil, fmt.Errorf("AssumeRole 请求失败: %w", err)
	}

	// 返回临时凭证
	logger.Log.Info("成功获取临时凭证")
	logger.Log.WithFields(logrus.Fields{
		"AccessKeyId":  *result.Credentials.AccessKeyId,
		"SessionToken": *result.Credentials.SessionToken,
		"Expiration":   *result.Credentials.Expiration,
	}).Debug("临时凭证详情")

	return &Credential{
		SecretId:  *result.Credentials.AccessKeyId,
		SecretKey: *result.Credentials.SecretAccessKey,
		Token:     *result.Credentials.SessionToken,
	}, nil
}
