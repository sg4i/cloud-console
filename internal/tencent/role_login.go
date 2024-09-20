package tencent

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/sg4i/cloud-console/internal/logger"
)

func GenerateRoleLoginSignature(secretId string, secretKey string, token string, nonce int, timestamp int64) string {
	param := map[string]string{
		"nonce":     fmt.Sprintf("%d", nonce),
		"timestamp": fmt.Sprintf("%d", timestamp),
		"secretId":  secretId,
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
	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(signStrTrimmed))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}

func GenerateRoleLoginURL(secretId string, secretKey string, token string, nonce int, timestamp int64, algorithm string, sUrl string) (string, error) {
	// 参数校验
	if secretId == "" {
		logger.Log.Error("secretId 不能为空")
		return "", errors.New("secretId 不能为空")
	}
	if secretKey == "" {
		logger.Log.Error("secretKey 不能为空")
		return "", errors.New("secretKey 不能为空")
	}
	if token == "" {
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

	signature := GenerateRoleLoginSignature(secretId, secretKey, token, nonce, timestamp)

	params := map[string]string{
		"algorithm": algorithm,
		"secretId":  secretId,
		"token":     token,
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
	logger.Log.WithField("url", generatedURL).Debug("生成的角色登录 URL")
	return generatedURL, nil
}
