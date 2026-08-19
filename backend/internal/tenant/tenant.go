package tenant

import (
	"context"
)

type ctxKey string

const tenantKey ctxKey = "tenant_id"

// WithTenant returns a context carrying the given tenant id.
func WithTenant(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, tenantKey, id)
}

// GetTenantID returns the tenant id from the context. When no tenant is set
// (e.g. before real auth is wired in) it falls back to 1 as a mock default.
func GetTenantID(ctx context.Context) int64 {
	if ctx != nil {
		if v, ok := ctx.Value(tenantKey).(int64); ok && v > 0 {
			return v
		}
	}
	return 1
}
