package tenant

import (
	"context"
	"testing"
)

func TestGetOrganizationID(t *testing.T) {
	// 1. Nil context
	if got := GetOrganizationID(nil); got != 1 {
		t.Fatalf("GetOrganizationID(nil) = %d, want 1", got)
	}

	// 2. Empty context
	if got := GetOrganizationID(context.Background()); got != 1 {
		t.Fatalf("GetOrganizationID(empty) = %d, want 1", got)
	}

	// 3. With 0 (invalid/fallback)
	if got := GetOrganizationID(WithTenant(context.Background(), 0)); got != 1 {
		t.Fatalf("GetOrganizationID() with 0 tenant = %d, want 1", got)
	}

	// 4. With typed OrganizationKey
	if got := GetOrganizationID(WithTenant(context.Background(), 2)); got != 2 {
		t.Fatalf("GetOrganizationID() = %d, want 2", got)
	}

	// 5. With string key fallback
	ctxStringKey := context.WithValue(context.Background(), string(OrganizationKey), int64(3))
	if got := GetOrganizationID(ctxStringKey); got != 3 {
		t.Fatalf("GetOrganizationID() with string key = %d, want 3", got)
	}
}
