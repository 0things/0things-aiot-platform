package tenant

import (
	"context"
)

type ctxKey string

const OrganizationKey ctxKey = "organization_id"

// WithOrganization stores the Logto organization identifier in the request context.
func WithOrganization(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, OrganizationKey, id)
}

// GetOrganizationID returns the Logto organization identifier. An absent
// organization is returned as an empty string and is never replaced by a default.
func GetOrganizationID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	// Gin Context exposes values stored with Set through its string-key lookup.
	if value, _ := ctx.Value(string(OrganizationKey)).(string); value != "" {
		return value
	}
	value, _ := ctx.Value(OrganizationKey).(string)
	return value
}
