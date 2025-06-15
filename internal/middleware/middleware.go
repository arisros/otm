package middleware

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

var (
	visitors     = map[string]*rate.Limiter{}
	cooldowns    = map[string]time.Time{}
	rateLimit    = rate.Every(1 * time.Second)
	burstLimit   = 3
	cooldownTime = 30 * time.Second
)

func limiterFor(ip string) *rate.Limiter {
	if l, ok := visitors[ip]; ok {
		return l
	}
	lim := rate.NewLimiter(rateLimit, burstLimit)
	visitors[ip] = lim
	return lim
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetIP(r)

		if until, blocked := cooldowns[ip]; blocked {
			if time.Now().Before(until) {
				time.Sleep(until.Sub(time.Now()))
				http.Error(w, "Too many requests 🙂", http.StatusTooManyRequests)
				return
			}
			delete(cooldowns, ip)
		}

		if !limiterFor(ip).Allow() {
			cooldowns[ip] = time.Now().Add(cooldownTime)
			log.Printf("⏱ %s temporarily throttled", ip)
			http.Error(w, "Too many requests 👍🏼", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	} else {
		ip = strings.Split(ip, ",")[0]
	}
	return strings.TrimSpace(ip)
}

func LookupCountry(ip string) string {
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://ip-api.com/json/" + ip + "?fields=country")
	if err != nil || resp.StatusCode != http.StatusOK {
		return "Unknown"
	}
	defer resp.Body.Close()

	var res struct {
		Country string `json:"country"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "Unknown"
	}
	return res.Country
}
