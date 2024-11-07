package signer

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"time"
)

// SignatureTool 提供签名和校验签名的功能
type Md5Signer struct {
	secretKey string
}

// NewSignatureTool 创建一个新的SignatureTool实例
func NewMd5Signer(secretKey string) *Md5Signer {
	return &Md5Signer{
		secretKey: secretKey,
	}
}

// GenerateSignature 生成签名
func (s *Md5Signer) GenerateSignature(params map[string]string, expireAt int64) (string, error) {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var signStr string
	for _, k := range keys {
		signStr += fmt.Sprintf("%s=%s&", k, params[k])
	}
	signStr += fmt.Sprintf("expireAt=%d&", expireAt)
	signStr += fmt.Sprintf("key=%s", s.secretKey)

	hasher := md5.New()
	_, err := hasher.Write([]byte(signStr)) // 将字符串写入hasher
	return hex.EncodeToString(hasher.Sum(nil)), err
}
func (s *Md5Signer) Method() SignerMethod {
	return SignerMd5
}

// VerifySignature 校验签名
func (s *Md5Signer) VerifySignature(params map[string]string, signature string, expireAt int64) (bool, error) {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signStr string
	for _, k := range keys {
		signStr += fmt.Sprintf("%s=%s&", k, params[k])
	}
	signStr += fmt.Sprintf("expireAt=%d&", expireAt)
	signStr += fmt.Sprintf("key=%s", s.secretKey)
	hasher := md5.New()
	hasher.Write([]byte(signStr)) // 将字符串写入hasher
	expectedSignature := hex.EncodeToString(hasher.Sum(nil))

	if expectedSignature != signature {
		return false, errors.New("签名不正确")
	}
	if time.Now().Unix() > expireAt {
		return false, errors.New("签名已过期")
	}
	return true, nil
}
