//go:build windows
// +build windows

package ipckit

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	gw "github.com/Microsoft/go-winio"
)

func (s *UnixDomainService) Serve() (chdata chan []byte, cherr chan error) {
	udsPath := s.udspath()
	chdata = make(chan []byte, 10)
	cherr = make(chan error, 2)
	if _, err := os.Stat(udsPath); err == nil {
		if err := os.RemoveAll(udsPath); err != nil {
			cherr <- err
			return
		}
	}
	// 启动 UDS 服务端（网络类型：unix）
	listener, err := gw.ListenPipe(udsPath, &gw.PipeConfig{
		// 管道模式：消息模式（MessageMode=true）或字节流模式（默认 false）
		// 消息模式下，Read 会按发送的消息边界返回数据，不会粘包
		MessageMode: true,
		// 输入缓冲区大小（默认 4096 字节）
		InputBufferSize: int32(s.BufferSizeInput),
		// 输出缓冲区大小（默认 4096 字节）
		OutputBufferSize: int32(s.BufferSizeOut),
		// 权限控制（可选）：限制哪些用户可访问管道
		// 示例：允许当前用户完全控制
		SecurityDescriptor: s.SecurityDescriptor,
	})
	if err != nil {
		cherr <- err
		return
	}
	defer func() {
		os.Remove(udsPath) // 退出时清理 .sock 文件
	}()
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
	conn, err := gw.DialPipe(udsPath, nil)
	if err != nil {
		return 0, fmt.Errorf("客户端连接失败：%w", err)
	}
	defer conn.Close()
	return conn.Write([]byte(input + "\n"))
}

func (s *UnixDomainService) WriteBytes(input []byte) (num int, err error) {
	udsPath := s.udspath()
	// 连接 UDS 服务端
	conn, err := gw.DialPipe(udsPath, nil)
	if err != nil {
		return 0, fmt.Errorf("客户端连接失败：%w", err)
	}
	defer conn.Close()
	return conn.Write(input)
}
