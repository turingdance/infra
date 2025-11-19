package netkit

import (
	"net"
	"net/http"
	"strings"
)

func RandomPort() (port int, err error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, err
}

// GetRealIP 获取客户端真实 IP（兼容代理转发场景）
func RealIP(r *http.Request) string {
	// 1. 优先从 X-Forwarded-For 中获取（多代理场景）
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For 格式：clientIP, proxy1IP, proxy2IP...
		ips := strings.Split(xForwardedFor, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			// 过滤空IP和内网IP，返回第一个有效公网IP
			if ip != "" && !IsPrivateIP(ip) {
				return ip
			}
		}
		// 若所有IP都是内网，取第一个（可能是局域网内代理场景）
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 2. 从 X-Real-IP 中获取（单代理场景，如 Nginx 直接转发）
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" && !IsPrivateIP(xRealIP) {
		return xRealIP
	}

	// 3. 从 RemoteAddr 中获取（直连场景，无代理）
	remoteAddr := r.RemoteAddr
	if remoteAddr != "" {
		// RemoteAddr 格式：ip:port，需剥离端口
		ip, _, err := net.SplitHostPort(remoteAddr)
		if err == nil && ip != "" && !IsPrivateIP(ip) {
			return ip
		}
		return ip // 直连内网场景（如本地开发）
	}

	// 4. 兜底：无有效IP时返回空字符串
	return ""
}

// isPrivateIP 判断IP是否为内网IP（避免被伪造请求头欺骗）
func IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return true // 无效IP视为内网（防止恶意输入）
	}

	// 私有IP网段：
	// 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16, 127.0.0.0/8（本地回环）
	privateNets := []*net.IPNet{
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
		{IP: net.IPv4(127, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		{IP: net.IPv6loopback, Mask: net.CIDRMask(128, 128)}, // IPv6 本地回环
	}

	for _, net := range privateNets {
		if net.Contains(ip) {
			return true
		}
	}
	return false
}
