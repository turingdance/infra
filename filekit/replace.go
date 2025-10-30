package filekit

import (
	"os"
	"strings"
)

// 把filepath 文件中的  searchStr  替换成  replaceStr
func ReplaceInFile(filePath, searchStr, replaceStr string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	newContent := strings.ReplaceAll(string(content), searchStr, replaceStr)
	return os.WriteFile(filePath, []byte(newContent), 0644)
}
