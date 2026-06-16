package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

var DefaultCORSConfig = CORSConfig{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id"},
	AllowCredentials: true,
	MaxAge:           86400,
}

func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	origins := cfg.AllowedOrigins
	if len(origins) == 0 {
		origins = DefaultCORSConfig.AllowedOrigins
	}

	methods := cfg.AllowedMethods
	if len(methods) == 0 {
		methods = DefaultCORSConfig.AllowedMethods
	}

	headers := cfg.AllowedHeaders
	if len(headers) == 0 {
		headers = DefaultCORSConfig.AllowedHeaders
	}

	allowAllOrigins := slices.Contains(origins, "*")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if allowAllOrigins {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if origin != "" && matchOrigin(origin, origins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))

			if cfg.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if cfg.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", cfg.MaxAge))
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func matchOrigin(origin string, allowed []string) bool {
	for _, a := range allowed {
		if a == origin {
			return true
		}
		if strings.HasPrefix(a, "*.") && strings.HasSuffix(origin, a[1:]) {
			return true
		}
	}
	return false
}
