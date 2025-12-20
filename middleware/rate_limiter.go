package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	mu       sync.Mutex
	lastSeen = map[string]time.Time{}
)

func RateLimit(next http.Handler) http.Handler {
	log.Println("RATE LIMIT MIDDLEWARE HIT")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			ip = host
		}

		now := time.Now()

		mu.Lock()
		last, ok := lastSeen[ip]
		if ok && now.Sub(last) < 2*time.Second {
			mu.Unlock()
			http.Error(w, "rate limit: 1 request every 2sec", http.StatusTooManyRequests)
			return
		}
		lastSeen[ip] = now
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
