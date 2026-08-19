package tenant

import (
	"context"
	"testing"
)

func TestGetTenantID(t *testing.T) {
	if got := GetTenantID(context.Background()); got != 1 {
		t.Fatalf("GetTenantID() = %d, want 1", got)
	}
	if got := GetTenantID(WithTenant(context.Background(), 0)); got != 1 {
		t.Fatalf("GetTenantID() with legacy tenant = %d, want 1", got)
	}
	if got := GetTenantID(WithTenant(context.Background(), 2)); got != 2 {
		t.Fatalf("GetTenantID() = %d, want 2", got)
	}
}
