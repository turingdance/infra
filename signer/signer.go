package signer

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SignerMethod string

const (
	SignerMd5   SignerMethod = "md5"
	SignerSh256 SignerMethod = "sh256"
)

type ISigner interface {
	GenerateSignature(params map[string]string, expireAt int64) (string, error)
	VerifySignature(params map[string]string, signature string, expireAt int64) (bool, error)
	Method() SignerMethod
}

type Signer struct {
	secretKey   string
	expireField string
	signField   string
	method      SignerMethod
}

type SignerOption func(*Signer)

func UseExpireAtField(f string) SignerOption {
	return func(ss *Signer) {
		ss.expireField = f
	}
}

func UseSignField(f string) SignerOption {
	return func(ss *Signer) {
		ss.signField = f
	}
}

func UseMethod(m SignerMethod) SignerOption {
	return func(ss *Signer) {
		ss.method = m
	}
}

// NewSignatureTool 创建一个新的SignatureTool实例
func New(method SignerMethod, secretKey string, options ...SignerOption) *Signer {
	r := &Signer{
		secretKey:   secretKey,
		signField:   "sign",
		expireField: "expireAt",
		method:      method,
	}
	for _, opt := range options {
		opt(r)
	}
	return r
}

func (s *Signer) Sign(params map[string]string, duration time.Duration) (signature string, err error) {
	var keys []string
	params[s.expireField] = fmt.Sprintf("%d", time.Now().Add(duration).Unix())
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var signStr string
	for _, k := range keys {
		// 不参与签名
		if k != s.signField {
			signStr += fmt.Sprintf("%s=%s&", k, params[k])
		}
	}
	signStr = strings.TrimSuffix(signStr, "&")
	//
	if s.method == SignerMd5 {
		hasher := md5.New()
		signStr += fmt.Sprintf("key=%s", s.secretKey)
		_, err = hasher.Write([]byte(signStr)) // 将字符串写入hasher
		return hex.EncodeToString(hasher.Sum(nil)), err
	}
	if s.method == SignerSh256 {
		mac := hmac.New(sha256.New, []byte(s.secretKey))
		mac.Write([]byte(signStr))
		signature = hex.EncodeToString(mac.Sum(nil))
	}
	return signature, err

}

// VerifySignature 校验签名
func (s *Signer) Verify(params map[string]string, signature string) (bool, error) {
	var keys []string
	for k := range params {
		if k != s.signField {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var signStr string
	for _, k := range keys {
		signStr += fmt.Sprintf("%s=%s&", k, params[k])
	}
	signStr = strings.TrimSuffix(signStr, "&")
	expireAt, err := strconv.Atoi(params[s.expireField])
	if err != nil {
		return false, err
	}
	//
	if int64(expireAt) > time.Now().Unix() {
		return false, errors.New("签名已过期")
	}
	//
	if s.method == SignerMd5 {
		hasher := md5.New()
		hasher.Write([]byte(signStr)) // 将字符串写入hasher
		expectedSignature := hex.EncodeToString(hasher.Sum(nil))
		if expectedSignature != signature {
			return false, errors.New("签名不正确")
		}
	}
	//
	if s.method == SignerSh256 {
		mac := hmac.New(sha256.New, []byte(s.secretKey))
		mac.Write([]byte(signStr))
		expectedSignature := hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
			return false, errors.New("签名不正确")
		}
	}
	return true, nil
}
