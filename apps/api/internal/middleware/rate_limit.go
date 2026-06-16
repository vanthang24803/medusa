package middleware

import (
	"net/http"

	"go.uber.org/ratelimit"
)

func RateLimit(rps int) func(http.Handler) http.Handler {
	rl := ratelimit.New(rps)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rl.Take()
			next.ServeHTTP(w, r)
		})
	}
}
