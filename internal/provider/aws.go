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

	return sts.New(sess), nil
}

func AwsAssumeRole(accessKeyId, secretAccessKey, sessionToken, roleArn string) (*Credential, error) {
	opts := &AssumeRoleOptions{
		RoleSessionName: DefaultRoleSessionName,
		DurationSeconds: DefaultDurationSeconds,
		RoleArn:         roleArn,
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
	logger.Log.Infof("调用AssumeRole成功获取角色%s的临时凭证", opts.RoleArn)
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
