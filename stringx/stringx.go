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
