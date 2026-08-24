package tenant

import (
	"context"
)

type ctxKey string

const organizationKey ctxKey = "organization_id"

// WithTenant returns a context carrying the given tenant id.
func WithTenant(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, organizationKey, id)
}

// GetOrganizationID returns the tenant id from the context. When no tenant is set
// (e.g. before real auth is wired in) it falls back to 1 as a mock default.
func GetOrganizationID(ctx context.Context) int64 {
	if ctx != nil {
		if v, ok := ctx.Value(organizationKey).(int64); ok && v > 0 {
			return v
		}
	}
	return 1
}
