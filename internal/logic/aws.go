package logic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sg4i/cloud-console/internal/provider"
	"github.com/sirupsen/logrus"
)

// AwsLoginOptions 定义AWS登录所需的参数
type AwsLoginOptions struct {
	Credential        *provider.Credential
	LoginURL          string // 登录成功后的跳转地址
	Destination       string // 目标URL，默认为AWS控制台
	RoleArn           string // Role ARN
	AssumeRoleOptions *provider.AssumeRoleOptions
}

const (
	// awsDefaultDestination 定义默认的目标URL
	awsDefaultDestination = "https://console.aws.amazon.com"
	// awsDefaultLoginURL 定义默认的登录URL
	awsDefaultLoginURL = "https://signin.aws.amazon.com/federation"
)

// awsSigninTokenResponse 定义获取签名令牌的响应结构
type awsSigninTokenResponse struct {
	SigninToken string `json:"SigninToken"`
}

// getAwsSigninToken 获取AWS签名令牌
func getAwsSigninToken(credential *provider.Credential) (string, error) {
	baseURL := awsDefaultLoginURL

	// 格式化临时凭证为JSON
	urlCredentials := map[string]string{
		"sessionId":    credential.SecretId,
		"sessionKey":   credential.SecretKey,
		"sessionToken": credential.Token,
	}
	jsonStringWithTempCredentials, err := json.Marshal(urlCredentials)
	if err != nil {
		return "", fmt.Errorf("临时凭证JSON序列化失败: %v", err)
	}

	// 请求AWS联合端点获取登录令牌
	formData := url.Values{}
	formData.Set("Action", "getSigninToken")
	formData.Set("SessionDuration", "43200") // 12小时会话时长
	formData.Set("Session", string(jsonStringWithTempCredentials))

	resp, err := http.PostForm(baseURL, formData)
	if err != nil {
		return "", fmt.Errorf("请求签名令牌失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取签名令牌失败, 状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var tokenResp awsSigninTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("解析响应JSON失败: %v", err)
	}

	logger.Log.WithFields(logrus.Fields{
		"SigninToken": tokenResp.SigninToken,
	}).Debug("获取签名令牌响应")

	return tokenResp.SigninToken, nil
}

// GenerateAwsRoleLoginURL 生成AWS角色登录链接
func GenerateAwsRoleLoginURL(opts *AwsLoginOptions) (string, error) {
	// 如果没有有效的 token，则使用 AssumeRole 获取临时密钥
	if opts.Credential.Token == "" {
		tempCred, err := provider.AwsAssumeRole(opts.Credential.SecretId, opts.Credential.SecretKey, "", opts.AssumeRoleOptions)
		if err != nil {
			return "", fmt.Errorf("获取临时密钥失败: %w", err)
		}
		opts.Credential = tempCred
	}

	// 获取 SigninToken
	signinToken, err := getAwsSigninToken(opts.Credential)
	if err != nil {
		return "", fmt.Errorf("获取SigninToken失败: %v", err)
	}

	// 如果 Destination 为空，使用默认值
	if opts.Destination == "" {
		opts.Destination = awsDefaultDestination
	}

	// 如果 LoginURL 为空，使用默认值
	if opts.LoginURL == "" {
		opts.LoginURL = awsDefaultLoginURL
	}

	queryParams := url.Values{}
	queryParams.Add("Action", "login")
	queryParams.Add("Issuer", opts.LoginURL)
	queryParams.Add("Destination", opts.Destination)
	queryParams.Add("SigninToken", signinToken)

	loginURL := fmt.Sprintf("%s?%s", awsDefaultLoginURL, queryParams.Encode())
	logger.Log.WithField("url", loginURL).Debug("生成的AWS角色登录 URL")

	return loginURL, nil
}
