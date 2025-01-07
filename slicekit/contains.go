package slicekit

import "strings"

func Contains[T comparable](arr []T, ele T) bool {
	result := false
	for _, v := range arr {
		if v == ele {
			result = true
			break
		}
	}
	return result
}

func HasElement[T comparable](arr []T, str T) bool {
	return Contains(arr, str)
}
func HasSubStr(strarr []string, str string) bool {
	return Contains(strarr, str)
}
func HasSubStrIgnoreCase(arr []string, ele string) bool {
	result := false
	for _, str := range arr {
		if strings.EqualFold(str, ele) {
			result = true
			break
		}
	}
	return result
}
