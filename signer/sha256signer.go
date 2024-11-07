package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// SignatureTool 提供签名和校验签名的功能
type Sha256Signer struct {
	secretKey string
}

// NewSignatureTool 创建一个新的SignatureTool实例
func NewSha256Signer(secretKey string) *Sha256Signer {
	return &Sha256Signer{
		secretKey: secretKey,
	}
}
func (s *Sha256Signer) Method() SignerMethod {
	return SignerSh256
}

// GenerateSignature 生成签名
func (s *Sha256Signer) GenerateSignature(params map[string]string, expireAt int64) (string, error) {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var signStr string
	for _, k := range keys {
		signStr += fmt.Sprintf("%s=%s&", k, params[k])
	}
	signStr = strings.TrimSuffix(signStr, "&")

	mac := hmac.New(sha256.New, []byte(s.secretKey))
	mac.Write([]byte(signStr))
	mac.Write([]byte(fmt.Sprintf("&expireAt=%d", expireAt)))
	signature := hex.EncodeToString(mac.Sum(nil))
	return signature, nil
}

// VerifySignature 校验签名
func (s *Sha256Signer) VerifySignature(params map[string]string, signature string, expireAt int64) (bool, error) {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signStr string
	for _, k := range keys {
		signStr += fmt.Sprintf("%s=%s&", k, params[k])
	}
	signStr = strings.TrimSuffix(signStr, "&")

	mac := hmac.New(sha256.New, []byte(s.secretKey))
	mac.Write([]byte(signStr))
	mac.Write([]byte(fmt.Sprintf("&expireAt=%d", expireAt)))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return false, errors.New("签名不正确")
	}
	if time.Now().Unix() > expireAt {
		return false, errors.New("签名已过期")
	}
	return true, nil
}
