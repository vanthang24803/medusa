package actor

import "context"

type key struct{}

// Set stores the current actor ID (auth identity ID or "apikey:<id>") in ctx.
func Set(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key{}, id)
}

// Get returns a pointer to the actor ID, or nil if not set.
func Get(ctx context.Context) *string {
	if v, ok := ctx.Value(key{}).(string); ok && v != "" {
		return &v
	}
	return nil
}
