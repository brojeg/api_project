package model

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
)

type HTTPReqInfo struct {
	method    string
	uri       string
	ipaddr    string
	code      int
	duration  time.Duration
	userAgent string
}

var HTTPLogger = func(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ri := &HTTPReqInfo{
			method:    r.Method,
			uri:       r.URL.String(),
			userAgent: r.Header.Get("User-Agent"),
		}

		ri.ipaddr = requestGetRemoteAddress(r)
		m := httpsnoop.CaptureMetrics(h, w, r)

		ri.code = m.Code
		ri.duration = m.Duration
		fmt.Println(ri) //TODO create a better representation
	}
	return http.HandlerFunc(fn)
}

func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		return parts[0]
	}
	return hdrRealIP
}
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}
