package middleware

import (
	"context"
	"net/http"

	"ecommerce/packages/httpx"
	"ecommerce/packages/types"
)

// iamEvaluator is a local interface to avoid importing the iam package.
type iamEvaluator interface {
	Evaluate(ctx context.Context, authIdentityID, action, resource string) (bool, error)
	EvaluateAPIKey(ctx context.Context, apiKeyID, action, resource string) (bool, error)
}

// RequirePermission enforces that the authenticated principal has the given action
// on the derived resource (service extracted from the action prefix, e.g. "product:Create" → "product/*").
func RequirePermission(iamSvc iamEvaluator, action string) func(http.Handler) http.Handler {
	resource := actionToResource(action)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if allowed := evaluate(r.Context(), iamSvc, action, resource); !allowed {
				httpx.Error(w, r, types.NewForbidden("insufficient permissions"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyPermission allows the request if the principal passes ANY of the given actions.
func RequireAnyPermission(iamSvc iamEvaluator, actions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, action := range actions {
				resource := actionToResource(action)
				if allowed := evaluate(r.Context(), iamSvc, action, resource); allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			httpx.Error(w, r, types.NewForbidden("insufficient permissions"))
		})
	}
}

// evaluate dispatches to the correct evaluator based on which principal context key is set.
func evaluate(ctx context.Context, iamSvc iamEvaluator, action, resource string) bool {
	if apiKeyID := GetAPIKeyID(ctx); apiKeyID != "" {
		ok, _ := iamSvc.EvaluateAPIKey(ctx, apiKeyID, action, resource)
		return ok
	}
	authID := GetAuthIdentityID(ctx)
	if authID == "" {
		return false
	}
	ok, _ := iamSvc.Evaluate(ctx, authID, action, resource)
	return ok
}

// actionToResource derives a service-level resource string from an action.
// "product:Create" → "product/*", "*" → "*"
func actionToResource(action string) string {
	if action == "*" {
		return "*"
	}
	for i, c := range action {
		if c == ':' {
			return action[:i] + "/*"
		}
	}
	return "*"
}
