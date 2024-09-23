package alibaba

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sg4i/cloud-console/internal/alibaba/config"
)

// RoleLoginParams 定义生成角色登录链接所需的参数
type RoleLoginParams struct {
	LoginURL    string
	Destination string
	Credential  config.Credential
}

// DefaultDestination 定义默认的目标URL
const DefaultDestination = "https://console.aliyun.com"

// GenerateRoleLoginURL 生成阿里云角色登录链接
func GenerateRoleLoginURL(params RoleLoginParams) (string, error) {
	// 获取 SigninToken
	signinToken, err := GetSigninToken(params.Credential)
	if err != nil {
		return "", fmt.Errorf("获取SigninToken失败: %v", err)
	}

	baseURL := "https://signin.aliyun.com/federation"

	// 如果 Destination 为空，使用默认值
	if params.Destination == "" {
		params.Destination = DefaultDestination
	}

	queryParams := url.Values{}
	queryParams.Add("Action", "Login")
	queryParams.Add("LoginUrl", params.LoginURL)
	queryParams.Add("Destination", params.Destination)
	queryParams.Add("SigninToken", signinToken)

	return fmt.Sprintf("%s?%s", baseURL, queryParams.Encode()), nil
}

// GetSigninTokenResponse 定义获取签名令牌的响应结构
type GetSigninTokenResponse struct {
	RequestId   string `json:"RequestId"`
	SigninToken string `json:"SigninToken"`
}

// GetSigninToken 获取阿里云签名令牌
func GetSigninToken(credential config.Credential) (string, error) {
	baseURL := "https://signin.aliyun.com/federation"

	formData := url.Values{}
	formData.Set("Action", "GetSigninToken")
	formData.Set("AccessKeyId", credential.AccessKeyId)
	formData.Set("AccessKeySecret", credential.AccessKeySecret)
	formData.Set("SecurityToken", credential.SecurityToken)
	formData.Set("TicketType", "mini")

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

	var tokenResp GetSigninTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("解析响应JSON失败: %v", err)
	}

	return tokenResp.SigninToken, nil
}
