//go:build darwin
// +build darwin

package ipckit

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
)

func (s *UnixDomainService) ServeString() (chdata chan string, cherr chan error) {
	udsPath := s.udspath()

	if _, err := os.Stat(udsPath); err == nil {
		if err := os.RemoveAll(udsPath); err != nil {
			cherr <- err
			return
		}
	}

	chdata = make(chan string, 1024)
	cherr = make(chan error, 10)
	// 启动 UDS 服务端（网络类型：unix）
	listener, err := net.Listen("unix", udsPath)
	if err != nil {
		cherr <- err
		return
	}
	defer func() {
		listener.Close()

		os.Remove(udsPath) // 退出时清理 .sock 文件

	}()
	// 并发处理客户端连接（每个连接开 1 个 goroutine，轻量无压力）
	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Write([]byte(err.Error()))
			continue
		}
		go s.handleClientConn(conn, chdata, cherr) // 异步处理，不阻塞主循环
	}
}

func (s *UnixDomainService) handleClientConn(conn net.Conn, chdata chan string, cherr chan error) {
	defer conn.Close() // 连接结束自动关闭

	// 带缓冲的读写器（提升小数据读写效率，减少系统调用）
	reader := bufio.NewReader(conn)
	for {
		// 读取客户端数据（阻塞，直到收到数据或连接断开）
		data, err := reader.ReadString('\n') // 按换行符分割数据（自定义分隔符也可）
		if err != nil {
			cherr <- err
			conn.Write([]byte(err.Error()))
			return
		}
		chdata <- data
	}
}

func (s *UnixDomainService) ServeBytes() (chdata chan []byte, cherr chan error) {
	udsPath := s.udspath()
	if runtime.GOOS != "windows" {
		_ = os.Remove(udsPath)
	}
	chdata = make(chan []byte, 1024)
	cherr = make(chan error, 10)
	// 启动 UDS 服务端（网络类型：unix）
	listener, err := net.Listen("unix", udsPath)
	if err != nil {
		cherr <- err
	}
	defer func() {
		_ = listener.Close()
		if runtime.GOOS != "windows" {
			_ = os.Remove(udsPath) // 退出时清理 .sock 文件
		}
	}()
	// 并发处理客户端连接（每个连接开 1 个 goroutine，轻量无压力）
	for {
		conn, err := listener.Accept()
		if err != nil {
			cherr <- err
			continue
		}
		go s.handleByteClientConn(conn, chdata, cherr) // 异步处理，不阻塞主循环
	}
}

func (s *UnixDomainService) handleByteClientConn(conn net.Conn, chdata chan []byte, cherr chan error) {
	defer conn.Close() // 连接结束自动关闭
	reader := bufio.NewReader(conn)
	for {
		// 读取客户端数据（阻塞，直到收到数据或连接断开）
		data := make([]byte, 0)
		size, err := reader.Read(data)
		if err != nil {
			cherr <- err
			return
		}
		chdata <- data[:size]
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
