package middleware

import (
	"context"
	"net/http"
	"strings"

	"ecommerce/packages/actor"
	"ecommerce/packages/httpx"
	"ecommerce/packages/types"
)

type contextKey string

const (
	CtxAuthIdentityID contextKey = "auth_identity_id"
	CtxCustomerID     contextKey = "customer_id"
	CtxAPIKeyID       contextKey = "api_key_id"
)

func GetAuthIdentityID(ctx context.Context) string {
	v, _ := ctx.Value(CtxAuthIdentityID).(string)
	return v
}

func GetCustomerID(ctx context.Context) string {
	v, _ := ctx.Value(CtxCustomerID).(string)
	return v
}

func GetAPIKeyID(ctx context.Context) string {
	v, _ := ctx.Value(CtxAPIKeyID).(string)
	return v
}

type tokenValidator interface {
	ValidateToken(ctx context.Context, tokenStr string) (string, string, error)
}

func RequireAuth(v tokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				httpx.Error(w, r, types.ErrUnauthorized)
				return
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				httpx.Error(w, r, types.ErrUnauthorized)
				return
			}

			authID, customerID, err := v.ValidateToken(r.Context(), parts[1])
			if err != nil {
				httpx.Error(w, r, types.ErrUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CtxAuthIdentityID, authID)
			ctx = context.WithValue(ctx, CtxCustomerID, customerID)
			ctx = actor.Set(ctx, authID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
