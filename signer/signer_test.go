package signer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/turingdance/infra/signer"
)

func TestMd5(t *testing.T) {
	// 示例使用
	secretKey := "your_secret_key"
	expireAt := time.Now().Add(24 * time.Hour).Unix() // 设置过期时间为24小时后

	// 创建SignatureTool实例
	signatureTool := signer.NewMd5Signer(secretKey)

	// 生成签名
	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	signature, err := signatureTool.GenerateSignature(params, expireAt)
	if err != nil {
		fmt.Println("Error generating signature:", err)
		return
	}
	fmt.Println("Generated Signature:", signature)

	// 校验签名
	isValid, err := signatureTool.VerifySignature(params, signature, expireAt)
	fmt.Println("Is the signature valid?", isValid)

}

func TestSha256(t *testing.T) {
	// 示例使用
	secretKey := "your_secret_key"
	expireAt := time.Now().Add(24 * time.Hour).Unix() // 设置过期时间为24小时后

	// 创建SignatureTool实例
	signatureTool := signer.NewSha256Signer(secretKey)

	// 生成签名
	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	signature, err := signatureTool.GenerateSignature(params, expireAt)
	if err != nil {
		fmt.Println("Error generating signature:", err)
		return
	}
	fmt.Println("Generated Signature:", signature)

	// 校验签名
	isValid, err := signatureTool.VerifySignature(params, signature, expireAt)
	fmt.Println("Is the signature valid?", isValid)

}

func Test(t *testing.T) {
	// 示例使用
	secretKey := "your_secret_key"
	// expireAt := time.Now().Add(24 * time.Hour).Unix() // 设置过期时间为24小时后

	// 创建SignatureTool实例
	signatureTool := signer.New(signer.SignerMd5, secretKey)

	// 生成签名
	params := map[string]string{
		"param1": "value1",
		"param2": "value2",
	}
	signature, err := signatureTool.Sign(params, 10*time.Second)
	if err != nil {
		fmt.Println("Error generating signature:", err)
		return
	}
	fmt.Println("Generated Signature:", signature)

	// 校验签名
	isValid, err := signatureTool.Verify(params, signature)
	fmt.Println("Is the signature valid?", isValid)

}
