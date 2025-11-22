package ipckit

import (
	"fmt"
	"runtime"
	"strings"
)

type UnixDomainService struct {
	Name               string
	Delim              string
	BufferSizeInput    int
	BufferSizeOut      int
	SecurityDescriptor string
}
type UdsOption func(*UnixDomainService)

func SetDelim(delim string) UdsOption {
	return func(uds *UnixDomainService) {
		uds.Delim = delim
	}
}
func SetName(name string) UdsOption {
	return func(uds *UnixDomainService) {
		uds.Name = name
	}
}
func SetBufferSizeInput(size int) UdsOption {
	return func(uds *UnixDomainService) {
		uds.BufferSizeInput = size
	}
}
func SetBufferSizeOut(size int) UdsOption {
	return func(uds *UnixDomainService) {
		uds.BufferSizeOut = size
	}
}
func SetSecurityDescriptor(str string) UdsOption {
	return func(uds *UnixDomainService) {
		uds.SecurityDescriptor = str
	}
}
func NewUnixDomainService(name string, opts ...UdsOption) *UnixDomainService {
	s := &UnixDomainService{
		Name:               strings.ToLower(name),
		Delim:              "\n",
		BufferSizeInput:    256,
		BufferSizeOut:      256,
		SecurityDescriptor: "D:P(A;;GA;;;IU)",
	}
	for _, v := range opts {
		v(s)
	}
	return s
}
func (s *UnixDomainService) Path() string {
	return s.udspath()
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
