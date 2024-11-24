package stringx

import (
	"regexp"
	"strings"
)

// 驼峰命名
// A_b_c =>aBC
func CamelLcFirst(input string) string {
	return Lcfirst(CamelUcFirst(input))
}
func CamelUcFirst(input string) string {
	p := regexp.MustCompile(`_\w{1}`)
	return p.ReplaceAllStringFunc(input, func(s string) string {
		return strings.ToUpper(s[1:])
	})
}

// 下划线转驼峰命名
func UnderlineToCamelCase(input string) string {
	strArr := strings.Split(input, "_")
	result := make([]string, 0)
	for _, v := range strArr {
		result = append(result, Ucfirst(v))
	}
	return Lcfirst(strings.Join(result, ""))
}

// 下划线转驼峰命名
func UnderlineToUperCamelCase(input string) string {
	strArr := strings.Split(input, "_")
	result := make([]string, 0)
	for _, v := range strArr {
		result = append(result, Ucfirst(v))
	}
	return Ucfirst(strings.Join(result, ""))
}
