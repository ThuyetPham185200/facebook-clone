package utils

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ===== Helpers =====
// helper: thay {param} bằng value thực
func ReplaceParam(path, param, value string) string {
	return strings.ReplaceAll(path, "{"+param+"}", value)
}

func CopySafeHeaders(src, dst http.Header) {
	if ct := src.Get("Content-Type"); ct != "" {
		dst.Set("Content-Type", ct)
	}
	if acc := src.Get("Accept"); acc != "" {
		dst.Set("Accept", acc)
	}
}

func WritePlainError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(msg))
}

func NewRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
