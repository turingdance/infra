package pkkit

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func UseDeviceInfo(deviceInfo ...string) string {
	hash := md5.New()
	hash.Write([]byte(strings.Join(deviceInfo, "|")))
	return hex.EncodeToString(hash.Sum(nil))
}
