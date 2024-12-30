package logic

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/sg4i/cloud-console/internal/logger"

	"github.com/sg4i/cloud-console/internal/provider"
)

func generateRoleLoginSignature(credential *provider.Credential, nonce int, timestamp int64) string {
	param := map[string]string{
		"nonce":     fmt.Sprintf("%d", nonce),
		"timestamp": fmt.Sprintf("%d", timestamp),
		"secretId":  credential.SecretId,
		"action":    "roleLogin",
	}

	// 对参数进行排序
	keys := make([]string, 0, len(param))
	for k := range param {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建签名字符串
	var signStr strings.Builder
	signStr.WriteString("GETcloud.tencent.com/login/roleAccessCallback?")
	for _, k := range keys {
		signStr.WriteString(fmt.Sprintf("%s=%s&", k, url.QueryEscape(param[k])))
	}
	signStrTrimmed := strings.TrimSuffix(signStr.String(), "&")

	// 计算HMAC-SHA1签名
	h := hmac.New(sha1.New, []byte(credential.SecretKey))
	h.Write([]byte(signStrTrimmed))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}

func generateRoleLoginURL(credential *provider.Credential, nonce int, timestamp int64, algorithm string, sUrl string) (string, error) {
	// 参数校验
	if credential.SecretId == "" {
		logger.Log.Error("secretId 不能为空")
		return "", errors.New("secretId 不能为空")
	}
	if credential.SecretKey == "" {
		logger.Log.Error("secretKey 不能为空")
		return "", errors.New("secretKey 不能为空")
	}
	if credential.Token == "" {
		return "", errors.New("token 不能为空")
	}
	if nonce <= 0 {
		return "", fmt.Errorf("nonce 必须大于 0，当前值：%d", nonce)
	}
	if timestamp <= 0 {
		return "", fmt.Errorf("timestamp 必须大于 0，当前值：%d", timestamp)
	}
	if algorithm == "" {
		return "", errors.New("algorithm 不能为空")
	}
	if sUrl == "" {
		return "", errors.New("sUrl 不能为空")
	}

	signature := generateRoleLoginSignature(credential, nonce, timestamp)

	params := map[string]string{
		"algorithm": algorithm,
		"secretId":  credential.SecretId,
		"token":     credential.Token,
		"nonce":     fmt.Sprintf("%d", nonce),
		"timestamp": fmt.Sprintf("%d", timestamp),
		"signature": signature,
		"s_url":     sUrl,
	}

	baseURL := "https://cloud.tencent.com/login/roleAccessCallback"
	var queryParams []string
	for key, value := range params {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, url.QueryEscape(value)))
	}

	generatedURL := fmt.Sprintf("%s?%s", baseURL, strings.Join(queryParams, "&"))
	return generatedURL, nil
}

type TencentLoginOptions struct {
	Credential        *provider.Credential
	SUrl              string // SUrl
	RoleArn           string // Role ARN
	AssumeRoleOptions *provider.AssumeRoleOptions
}

func GenerateTencentRoleLoginURL(opts *TencentLoginOptions) (string, error) {
	// 如果没有有效的 token，则使用 AssumeRole 获取临时密钥
	if opts.Credential.Token == "" {
		tempCred, err := provider.TencentAssumeRole(opts.Credential.SecretId, opts.Credential.SecretKey, opts.RoleArn)
		if err != nil {
			return "", fmt.Errorf("获取临时密钥失败: %w", err)
		}
		opts.Credential.SecretId = tempCred.SecretId
		opts.Credential.SecretKey = tempCred.SecretKey
		opts.Credential.Token = tempCred.Token
	}

	algorithm := "sha1"
	nonce := rand.Intn(900000) + 100000
	timestamp := time.Now().Unix()

	url, err := generateRoleLoginURL(opts.Credential, nonce, timestamp, algorithm, opts.SUrl)
	if err != nil {
		logger.Log.WithError(err).Error("生成腾讯云角色登录 URL 失败")
		return "", fmt.Errorf("生成腾讯云角色登录 URL 失败: %w", err)
	}

	logger.Log.WithField("url", url).Debug("成功生成腾讯云角色登录 URL")
	return url, nil
}
