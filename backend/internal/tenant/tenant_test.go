package tenant

import (
	"context"
	"testing"
)

func TestGetOrganizationID(t *testing.T) {
	if got := GetOrganizationID(context.Background()); got != 1 {
		t.Fatalf("GetOrganizationID() = %d, want 1", got)
	}
	if got := GetOrganizationID(WithTenant(context.Background(), 0)); got != 1 {
		t.Fatalf("GetOrganizationID() with legacy tenant = %d, want 1", got)
	}
	if got := GetOrganizationID(WithTenant(context.Background(), 2)); got != 2 {
		t.Fatalf("GetOrganizationID() = %d, want 2", got)
	}
}
