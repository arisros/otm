package middleware

import (
	"net/http"
	"sync"
	"time"
)

type rateData struct {
	Count     int
	LastReset time.Time
}

var limiter = struct {
	sync.Mutex
	data map[string]*rateData
}{data: make(map[string]*rateData)}

const (
	LimitRequestCount = 10               // Max 10 requests
	LimitDuration     = 20 * time.Second // Per 20 seconds
)

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		now := time.Now()

		limiter.Lock()
		info, exists := limiter.data[ip]
		if !exists || now.Sub(info.LastReset) > LimitDuration {
			limiter.data[ip] = &rateData{Count: 1, LastReset: now}
			limiter.Unlock()
		} else {
			if info.Count >= LimitRequestCount {
				limiter.Unlock()
				http.Error(w, "🚫 Too Many Requests", http.StatusTooManyRequests)
				return
			}
			info.Count++
			limiter.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}
