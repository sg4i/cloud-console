package aws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sg4i/cloud-console/internal/aws/config"
	"github.com/sg4i/cloud-console/internal/logger"
	"github.com/sirupsen/logrus"
)

// RoleLoginParams 定义生成角色登录链接所需的参数
type RoleLoginParams struct {
	LoginURL    string
	Destination string
	Credential  config.Credential
}

// DefaultDestination 定义默认的目标URL
const DefaultDestination = "https://console.aws.amazon.com"

// DefaultLoginURL 定义默认的登录URL
const DefaultLoginURL = "https://signin.aws.amazon.com/federation"

// GenerateRoleLoginURL 生成AWS角色登录链接
func GenerateRoleLoginURL(params RoleLoginParams) (string, error) {
	// 获取 SigninToken
	signinToken, err := GetSigninToken(params.Credential)
	if err != nil {
		return "", fmt.Errorf("获取SigninToken失败: %v", err)
	}

	baseURL := "https://signin.aws.amazon.com/federation"

	// 如果 Destination 为空，使用默认值
	if params.Destination == "" {
		params.Destination = DefaultDestination
	}

	// 如果 LoginURL 为空，使用默认值
	if params.LoginURL == "" {
		params.LoginURL = DefaultLoginURL
	}

	queryParams := url.Values{}
	queryParams.Add("Action", "login")
	queryParams.Add("Issuer", params.LoginURL)
	queryParams.Add("Destination", params.Destination)
	queryParams.Add("SigninToken", signinToken)

	return fmt.Sprintf("%s?%s", baseURL, queryParams.Encode()), nil
}

// GetSigninTokenResponse 定义获取签名令牌的响应结构
type GetSigninTokenResponse struct {
	SigninToken string `json:"SigninToken"`
}

// GetSigninToken 获取AWS签名令牌
func GetSigninToken(credential config.Credential) (string, error) {
	baseURL := "https://signin.aws.amazon.com/federation"

	// Step 3: Format resulting temporary credentials into JSON
	urlCredentials := map[string]string{
		"sessionId":    credential.AccessKeyId,
		"sessionKey":   credential.SecretAccessKey,
		"sessionToken": credential.SessionToken,
	}
	jsonStringWithTempCredentials, err := json.Marshal(urlCredentials)
	if err != nil {
		return "", fmt.Errorf("临时凭证JSON序列化失败: %v", err)
	}

	// Step 4: Make request to AWS federation endpoint to get sign-in token
	formData := url.Values{}
	formData.Set("Action", "getSigninToken")
	formData.Set("SessionDuration", "43200") // a 12-hour session duration
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

	var tokenResp GetSigninTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("解析响应JSON失败: %v", err)
	}

	logger.Log.WithFields(logrus.Fields{
		"SigninToken": tokenResp.SigninToken,
	}).Debug("获取签名令牌响应")

	return tokenResp.SigninToken, nil
}
