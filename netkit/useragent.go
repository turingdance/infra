package netkit

import "strings"

// OSType 是一个字符串类型，用于表示操作系统
type OSType string

// 定义操作系统的字符串常量
const (
	OSUnknown   OSType = "Unknown"
	OSiOS       OSType = "iOS"
	OSAndroid   OSType = "Android"
	OSHarmonyOS OSType = "HarmonyOS"
	OSWindows   OSType = "Windows"
	OSMacOS     OSType = "macOS"
	OSLinux     OSType = "Linux"
)

// DeviceType 是一个字符串类型，用于表示设备类型
type DeviceType string

// 定义设备类型的字符串常量
const (
	DeviceUnknown DeviceType = "Unknown"
	DevicePC      DeviceType = "PC"
	DeviceMobile  DeviceType = "Mobile"
	DeviceTablet  DeviceType = "Tablet"
)

func (ua OSType) String() string {
	return string(ua)
}

// UAInfo 结构体存储了解析 User-Agent 后得到的核心信息
type UAInfo struct {
	Raw    string     // 原始的 User-Agent 字符串
	OS     OSType     // 操作系统类型 (string)
	Device DeviceType // 设备类型 (string)
}

// ParseUserAgent 解析 User-Agent 字符串并返回一个包含 OS 和 Device 的 UAInfo 对象
func ParseUserAgent(userAgent string) UAInfo {
	uaLower := strings.ToLower(userAgent)
	info := UAInfo{
		Raw: userAgent,
	}

	// --- 1. 解析操作系统 (OSType) ---
	switch {
	case strings.Contains(uaLower, "harmonyos"):
		info.OS = OSHarmonyOS
	case strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "ipod") || strings.Contains(uaLower, "ipad"):
		info.OS = OSiOS
	case strings.Contains(uaLower, "android"):
		info.OS = OSAndroid
	case strings.Contains(uaLower, "windows nt"):
		info.OS = OSWindows
	case strings.Contains(uaLower, "mac os x"):
		info.OS = OSMacOS
	case strings.Contains(uaLower, "linux"):
		// 排除掉已经识别为 Android 或 HarmonyOS 的情况
		if !strings.Contains(uaLower, "android") && !strings.Contains(uaLower, "harmonyos") {
			info.OS = OSLinux
		}
	default:
		info.OS = OSUnknown
	}

	// --- 2. 解析设备类型 (DeviceType) ---
	switch {
	case strings.Contains(uaLower, "ipad") || strings.Contains(uaLower, "tablet") || strings.Contains(uaLower, "pad"):
		info.Device = DeviceTablet
	case strings.Contains(uaLower, "mobile"):
		info.Device = DeviceMobile
	case info.OS == OSWindows || info.OS == OSMacOS || info.OS == OSLinux:
		// 桌面操作系统默认是 PC
		info.Device = DevicePC
	default:
		// 对于未知设备，根据已知的移动 OS 进行推断
		if info.OS == OSiOS || info.OS == OSAndroid || info.OS == OSHarmonyOS {
			info.Device = DeviceMobile
		} else {
			info.Device = DeviceUnknown
		}
	}

	return info
}
