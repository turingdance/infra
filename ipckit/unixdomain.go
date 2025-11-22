package ipckit

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"runtime"
	"strings"
)

type UnixDomainService struct {
	Name               string
	Patern             []byte
	BufferSizeInput    int
	BufferSizeOut      int
	SecurityDescriptor string
}
type UdsOption func(*UnixDomainService)

func SetDelim(delim string) UdsOption {
	return func(uds *UnixDomainService) {
		uds.Patern = []byte(delim)
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
		Patern:             []byte("\n"),
		BufferSizeInput:    256,
		BufferSizeOut:      256,
		SecurityDescriptor: "D:P(A;;GA;;;IU)",
	}
	for _, v := range opts {
		v(s)
	}
	return s
}
func Pipe(name string, opts ...UdsOption) *UnixDomainService {
	return NewUnixDomainService(name, opts...)
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

// 分割函数
func (s *UnixDomainService) Splitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// 1. 如果已经到达文件末尾且没有数据了，直接返回。
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	index := bytes.Index(data, s.Patern)
	if index == -1 {
		return 0, nil, nil
	}
	return index + len(s.Patern), data[:index], nil
}

func (s *UnixDomainService) handleByteClientConn(conn net.Conn, chdata chan []byte, cherr chan error) {
	defer conn.Close() // 连接结束自动关闭
	scanner := bufio.NewScanner(conn)
	scanner.Split(s.Splitter)
	for scanner.Scan() {
		pkg := scanner.Bytes()
		chdata <- pkg
	}
	if err := scanner.Err(); err != nil {
		if err == io.EOF {
			return
		} else {
			cherr <- err
		}
	}
}
