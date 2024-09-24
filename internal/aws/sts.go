package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/sg4i/cloud-console/internal/aws/config"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sirupsen/logrus"
)

// AssumeRoleOptions 定义 AssumeRole 的可选参数
type AssumeRoleOptions struct {
	RoleArn         string
	RoleSessionName string
	DurationSeconds int64
}

// DefaultAssumeRoleOptions 提供默认的 AssumeRole 选项
var DefaultAssumeRoleOptions = AssumeRoleOptions{
	RoleSessionName: "RoleSession",
	DurationSeconds: 3600,
}

func AssumeRole(accessKeyId, secretAccessKey, sessionToken string, opts *AssumeRoleOptions) (config.Credential, error) {
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
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyId, secretAccessKey, sessionToken),
		Region:      aws.String("us-east-1"),
	})
	if err != nil {
		logger.Log.WithError(err).Error("创建 STS 客户端失败")
		return config.Credential{}, fmt.Errorf("创建 STS 客户端失败: %w", err)
	}
	svc := sts.New(sess)

	// 创建 AssumeRole 请求
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(opts.RoleArn),
		RoleSessionName: aws.String(opts.RoleSessionName),
		DurationSeconds: aws.Int64(opts.DurationSeconds),
	}

	// 发送请求
	result, err := svc.AssumeRole(input)
	if err != nil {
		logger.Log.WithError(err).Error("AssumeRole 请求失败")
		return config.Credential{}, fmt.Errorf("AssumeRole 请求失败: %w", err)
	}

	// 返回临时凭证
	logger.Log.Info("成功获取临时凭证")
	logger.Log.WithFields(logrus.Fields{
		"AccessKeyId":  *result.Credentials.AccessKeyId,
		"SessionToken": *result.Credentials.SessionToken,
		"Expiration":   *result.Credentials.Expiration,
	}).Debug("临时凭证详情")
	return config.Credential{
		AccessKeyId:     *result.Credentials.AccessKeyId,
		SecretAccessKey: *result.Credentials.SecretAccessKey,
		SessionToken:    *result.Credentials.SessionToken,
	}, nil
}
