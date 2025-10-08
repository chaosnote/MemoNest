package utils

import (
	"encoding/json"
	"io"
	"net"
)

// IsPrivateIP 檢查 IP 地址是否為私有 IP（區域網路 IP）
func IsPrivateIP(ip net.IP) bool {
	// 檢查 IPv4 私有地址
	if ip.IsLoopback() || ip.IsPrivate() {
		return true
	}
	// 檢查 IPv6 私有地址 (Unique Local Address)
	if ip.To4() == nil && ip.IsGlobalUnicast() {
		return false
	}
	return false
}

func HttpRequestJSONUnmarshal(reader io.ReadCloser, output any) (e error) {
	var body []byte
	body, e = io.ReadAll(reader)
	if e != nil {
		return
	}
	defer reader.Close()

	e = json.Unmarshal(body, output)
	if e != nil {
		return
	}
	return
}
