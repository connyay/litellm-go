package ratelimit

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func NewRequestsPerMinuteLimiter(rpm int) func(http.Handler) http.Handler {
	if rpm <= 0 {
		// no limit
		return func(next http.Handler) http.Handler { return next }
	}
	lim := rate.NewLimiter(rate.Every(time.Minute/time.Duration(rpm)), rpm)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !lim.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte("rate limit exceeded"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
