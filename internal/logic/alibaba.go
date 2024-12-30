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

// AlibabaLoginOptions 定义阿里云登录所需的参数
type AlibabaLoginOptions struct {
	Credential        *provider.Credential
	LoginURL          string // 登录成功后的跳转地址
	Destination       string // 目标URL，默认为阿里云控制台
	RoleArn           string // Role ARN
	AssumeRoleOptions *provider.AssumeRoleOptions
}

// defaultDestination 定义默认的目标URL
const defaultDestination = "https://console.aliyun.com"

// getSigninTokenResponse 定义获取签名令牌的响应结构
type getSigninTokenResponse struct {
	RequestId   string `json:"RequestId"`
	SigninToken string `json:"SigninToken"`
}

// getSigninToken 获取阿里云签名令牌
func getSigninToken(credential *provider.Credential) (string, error) {
	baseURL := "https://signin.aliyun.com/federation"

	formData := url.Values{}
	formData.Set("Action", "GetSigninToken")
	formData.Set("AccessKeyId", credential.SecretId)
	formData.Set("AccessKeySecret", credential.SecretKey)
	formData.Set("SecurityToken", credential.Token)
	formData.Set("TicketType", "normal")

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

	var tokenResp getSigninTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("解析响应JSON失败: %v", err)
	}

	logger.Log.WithFields(logrus.Fields{
		"RequestId":   tokenResp.RequestId,
		"SigninToken": tokenResp.SigninToken,
	}).Debug("获取签名令牌响应")

	return tokenResp.SigninToken, nil
}

// GenerateAlibabaRoleLoginURL 生成阿里云角色登录链接
func GenerateAlibabaRoleLoginURL(opts *AlibabaLoginOptions) (string, error) {
	// 如果没有有效的 token，则使用 AssumeRole 获取临时密钥
	if opts.Credential.Token == "" {
		tempCred, err := provider.AlibabaAssumeRole(opts.Credential.SecretId, opts.Credential.SecretKey, opts.RoleArn)
		if err != nil {
			return "", fmt.Errorf("获取临时密钥失败: %w", err)
		}
		opts.Credential = tempCred
	}

	// 获取 SigninToken
	signinToken, err := getSigninToken(opts.Credential)
	if err != nil {
		return "", fmt.Errorf("获取SigninToken失败: %v", err)
	}

	baseURL := "https://signin.aliyun.com/federation"

	// 如果 Destination 为空，使用默认值
	if opts.Destination == "" {
		opts.Destination = defaultDestination
	}

	queryParams := url.Values{}
	queryParams.Add("Action", "Login")
	queryParams.Add("LoginUrl", opts.LoginURL)
	queryParams.Add("Destination", opts.Destination)
	queryParams.Add("SigninToken", signinToken)

	loginURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
	logger.Log.WithField("url", loginURL).Debug("生成的阿里云角色登录 URL")

	return loginURL, nil
}
