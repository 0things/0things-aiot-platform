package tenant

import (
	"context"
	"testing"
)

func TestGetOrganizationID(t *testing.T) {
	if got := GetOrganizationID(nil); got != "" {
		t.Fatalf("GetOrganizationID(nil) = %q, want empty", got)
	}

	if got := GetOrganizationID(context.Background()); got != "" {
		t.Fatalf("GetOrganizationID(empty) = %q, want empty", got)
	}

	if got := GetOrganizationID(WithOrganization(context.Background(), "org-1")); got != "org-1" {
		t.Fatalf("GetOrganizationID() = %q, want org-1", got)
	}

	ctxStringKey := context.WithValue(context.Background(), string(OrganizationKey), "org-3")
	if got := GetOrganizationID(ctxStringKey); got != "org-3" {
		t.Fatalf("GetOrganizationID() with string key = %q, want org-3", got)
	}
}
