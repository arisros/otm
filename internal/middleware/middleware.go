package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter)
var rateLimit = rate.Every(1 * time.Second)
var burstLimit = 3

// getVisitorLimiter returns a rate limiter for an IP address
func getVisitorLimiter(ip string) *rate.Limiter {
	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rateLimit, burstLimit)
		visitors[ip] = limiter
	}
	return limiter
}

// RateLimitMiddleware applies rate limiting per IP
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetIP(r)
		limiter := getVisitorLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetIP extracts the client IP address from headers or remote addr
func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	} else {
		ip = strings.Split(ip, ",")[0]
	}
	return strings.TrimSpace(ip)
}

// LookupCountry tries to identify the country for an IP address using ip-api.com
func LookupCountry(ip string) string {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://ip-api.com/json/" + ip + "?fields=country")
	if err != nil || resp.StatusCode != http.StatusOK {
		return "Unknown"
	}
	defer resp.Body.Close()

	var data struct {
		Country string `json:"country"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "Unknown"
	}
	return data.Country
}
