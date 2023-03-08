package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
)

type HTTPReqInfo struct {
	Method      string        `json:"method"`
	URI         string        `json:"uri"`
	IPAddr      string        `json:"ipaddr"`
	Code        int           `json:"code"`
	Duration    time.Duration `json:"duration"`
	UserAgent   string        `json:"user_agent"`
	RequestBody any
}

var HTTPLoggerMiddleware = func(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestInfo := &HTTPReqInfo{
			Method:    r.Method,
			URI:       r.URL.String(),
			UserAgent: r.Header.Get("User-Agent"),
		}

		requestInfo.IPAddr = requestGetRemoteAddress(r)
		m := httpsnoop.CaptureMetrics(h, w, r)

		requestInfo.Code = m.Code
		requestInfo.Duration = m.Duration
		data, err := json.Marshal(requestInfo)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data)) //TODO create a better representation
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
