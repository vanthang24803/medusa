package middleware

import (
	"context"
	"net/http"

	"ecommerce/modules/auth"
	"ecommerce/packages/actor"
	"ecommerce/packages/httpx"
	"ecommerce/packages/types"
)

const apiKeyHeader = "X-API-Key"

type apiKeyValidatorSvc interface {
	ValidateAPIKey(ctx context.Context, token string) (*auth.APIKey, error)
}

// RequireAPIKey authenticates requests via the X-API-Key header.
// On success, sets CtxAPIKeyID in the context for downstream permission checks.
func RequireAPIKey(v apiKeyValidatorSvc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(apiKeyHeader)
			if token == "" {
				httpx.Error(w, r, types.ErrUnauthorized)
				return
			}
			k, err := v.ValidateAPIKey(r.Context(), token)
			if err != nil {
				httpx.Error(w, r, types.ErrUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), CtxAPIKeyID, k.ID)
			ctx = actor.Set(ctx, "apikey:"+k.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAnyAuth accepts either a Bearer JWT or an X-API-Key header.
func RequireAnyAuth(jwt tokenValidator, apiKey apiKeyValidatorSvc) func(http.Handler) http.Handler {
	bearerMw := RequireAuth(jwt)
	apiKeyMw := RequireAPIKey(apiKey)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "" {
				bearerMw(next).ServeHTTP(w, r)
				return
			}
			if r.Header.Get(apiKeyHeader) != "" {
				apiKeyMw(next).ServeHTTP(w, r)
				return
			}
			httpx.Error(w, r, types.ErrUnauthorized)
		})
	}
}
