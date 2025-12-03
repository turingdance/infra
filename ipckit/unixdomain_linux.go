//go:build linux
// +build linux

package ipckit

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func (s *UnixDomainService) ServeBytes() (chdata chan []byte, cherr chan error) {
	udsPath := s.udspath()
	chdata = make(chan []byte, 1024)
	cherr = make(chan error, 10)
	// 启动 UDS 服务端（网络类型：unix）
	listener, err := net.Listen("unix", udsPath)
	if err != nil {
		cherr <- err
		return
	}
	// 并发处理客户端连接（每个连接开 1 个 goroutine，轻量无压力）
	go func() {
		for {
			if conn, err := listener.Accept(); err != nil {
				if err == io.EOF {

				} else {
					cherr <- err
					continue
				}
			} else {
				go s.handleByteClientConn(conn, chdata, cherr) // 异步处理，不阻塞主循环
			}

		}
	}()
	return
}

func (s *UnixDomainService) handleByteClientConn(conn net.Conn, chdata chan []byte, cherr chan error) {
	defer conn.Close() // 连接结束自动关闭
	reader := bufio.NewReader(conn)
	for {
		// 读取客户端数据（阻塞，直到收到数据或连接断开）
		data := make([]byte, s.BufferSizeInput)
		if size, err := reader.Read(data); err != nil {
			if err == io.EOF {
				return
			} else {
				cherr <- err
			}
		} else {
			chdata <- data[:size]
		}

	}
}

// white
func (s *UnixDomainService) WriteLine(input string) (num int, err error) {
	udsPath := s.udspath()
	// 连接 UDS 服务端
	conn, err := net.Dial("unix", udsPath)
	if err != nil {
		return 0, fmt.Errorf("客户端连接失败：%w", err)
	}
	defer conn.Close()
	return conn.Write([]byte(input + "\n"))
}

func (s *UnixDomainService) WriteString(input string) (num int, err error) {
	udsPath := s.udspath()
	// 连接 UDS 服务端
	conn, err := net.Dial("unix", udsPath)
	if err != nil {
		return 0, fmt.Errorf("客户端连接失败：%w", err)
	}
	defer conn.Close()
	return conn.Write([]byte(input))
}

func (s *UnixDomainService) WriteBytes(input []byte) (num int, err error) {
	udsPath := s.udspath()
	// 连接 UDS 服务端
	conn, err := net.Dial("unix", udsPath)
	if err != nil {
		return 0, fmt.Errorf("客户端连接失败：%w", err)
	}
	defer conn.Close()
	return conn.Write(input)
}
