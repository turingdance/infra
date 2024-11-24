package stringx

import (
	"html/template"
	"strings"
)

func JSStr(input string, args ...string) (result template.JSStr) {
	return template.JSStr(strings.Join(append([]string{input}, args...), ""))
}

func JS(input string, args ...string) (result template.JS) {
	return template.JS(strings.Join(append([]string{input}, args...), ""))
}

func Unescaped(str string) template.HTML {
	return template.HTML(str)
}
