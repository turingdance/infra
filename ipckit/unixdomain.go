package ipckit

import (
	"fmt"
	"runtime"
	"strings"
)

type UnixDomainService struct {
	Name string
}

func NewUnixDomainService(name string) *UnixDomainService {
	s := &UnixDomainService{
		Name: strings.ToLower(name),
	}

	return s
}
func (s *UnixDomainService) udspath() string {
	switch runtime.GOOS {
	case "windows":
		return `\\.\pipe\` + s.Name // Windows 固定前缀
	case "linux", "darwin": // Linux/macOS 统一用 /tmp 路径
		return fmt.Sprintf("/tmp/%s.sock", s.Name)
	default:
		return fmt.Sprintf("/tmp/%s.sock", s.Name)

	}
}
