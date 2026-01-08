package middleware

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Простая in-memory реализация rate limiter по IP (token bucket)
type bucket struct {
	tokens float64
	last   time.Time
	mu     sync.Mutex
}

var (
	buckets   = make(map[string]*bucket)
	bucketsMu sync.Mutex
	rps       = 5.0
	burst     = 10.0
)

func init() {
	if v := os.Getenv("RATE_LIMIT_RPS"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			rps = f
		}
	}
	if v := os.Getenv("RATE_LIMIT_BURST"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			burst = f
		}
	}
	// background cleanup
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			cutoff := time.Now().Add(-10 * time.Minute)
			bucketsMu.Lock()
			for k, b := range buckets {
				b.mu.Lock()
				last := b.last
				b.mu.Unlock()
				if last.Before(cutoff) {
					delete(buckets, k)
				}
			}
			bucketsMu.Unlock()
		}
	}()
}

func getBucket(key string) *bucket {
	bucketsMu.Lock()
	defer bucketsMu.Unlock()
	if b, ok := buckets[key]; ok {
		return b
	}
	b := &bucket{tokens: burst, last: time.Now()}
	buckets[key] = b
	return b
}

func allow(b *bucket) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now
	b.tokens += elapsed * rps
	if b.tokens > burst {
		b.tokens = burst
	}
	if b.tokens >= 1.0 {
		b.tokens -= 1.0
		return true
	}
	return false
}

// Middleware ограничивает количество запросов по IP
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		b := getBucket(ip)
		if !allow(b) {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
