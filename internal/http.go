package internal

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

// ExtractRequestClientIP ...
func ExtractRequestClientIP(req *http.Request) string {
	header := req.Header.Get("X-Real-IP")
	if header != "" {
		return header
	}
	header = req.Header.Get("X-Forwarded-For")
	if header != "" {
		return header
	}
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "unknown"
	}
	return host
}

// AddResponseRateLimitHeaders ...
func AddResponseRateLimitHeaders(resp http.ResponseWriter, limit int64, remaining int64, ttl time.Duration) {
	resp.Header().Add("X-Ratelimit-Limit", strconv.FormatInt(limit, 10))
	resp.Header().Add("X-Ratelimit-Remaining", strconv.FormatInt(max(remaining, 0), 10))
	resp.Header().Add("X-Ratelimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))
	if remaining < 0 {
		resp.Header().Add("Retry-After", strconv.FormatInt(int64(ttl.Seconds()), 10))
	}
}

func max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}
