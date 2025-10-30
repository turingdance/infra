package stringx

import (
	"strings"
)

func Ucfirst(input string) (result string) {
	if len(input) < 1 {
		return ""
	}
	return (strings.ToUpper(input[:1]) + input[1:])
}
func Upper(input string) (out string) {
	return strings.ToUpper(input)
}
func Lower(input string) (out string) {
	return strings.ToLower(input)
}
func Lcfirst(input string) (result string) {
	if len(input) < 1 {
		return ""
	}
	return (strings.ToLower(input[:1]) + input[1:])
}

// rightPad 右填充：在 s 右侧填充 padChar，使总长度为 length
func RightPad(s, padChar string, length int) string {
	if len(s) >= length {
		return s
	}
	padCount := length - len(s)
	return s + strings.Repeat(padChar, padCount)
}

// leftPad 左填充：在 s 左侧填充 padChar，使总长度为 length
func LeftPad(s, padChar string, length int) string {
	if len(s) >= length {
		return s // 原字符串长度已满足，直接返回
	}
	// 计算需要填充的数量
	padCount := length - len(s)
	// 生成填充字符串 + 原字符串
	return strings.Repeat(padChar, padCount) + s
}
