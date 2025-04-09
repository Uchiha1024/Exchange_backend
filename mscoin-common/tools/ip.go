package tools

import (
	"net"
	"net/http"
)

func GetRemoteClientIp(r *http.Request) string {
	remoteIp := r.RemoteAddr

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		remoteIp = ip
	} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
		remoteIp = ip
	} else {
		remoteIp, _, _ = net.SplitHostPort(remoteIp)
	}

	//本地ip
	if remoteIp == "::1" {
		remoteIp = "127.0.0.1"
	}

	return remoteIp
}
