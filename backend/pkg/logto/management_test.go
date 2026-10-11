package logto

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestManagementClientOrganizationAndInvitations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/oidc/token":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token": "mock-token",
				"expires_in":   3600,
			})
		case r.URL.Path == "/api/organization-invitations" && r.Method == http.MethodPost:
			var req CreateInvitationRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			_ = json.NewEncoder(w).Encode(OrganizationInvitation{
				ID:             "inv-1",
				InviterID:      req.InviterID,
				Invitee:        req.Invitee,
				OrganizationID: req.OrganizationID,
				Status:         "Pending",
				CreatedAt:      1700000000,
				ExpiresAt:      req.ExpiresAt,
			})
		case r.URL.Path == "/api/organization-invitations" && r.Method == http.MethodGet:
			orgID := r.URL.Query().Get("organizationId")
			require.Equal(t, "org-1", orgID)
			_ = json.NewEncoder(w).Encode([]OrganizationInvitation{
				{
					ID:             "inv-1",
					InviterID:      "user-1",
					Invitee:        "test@example.com",
					OrganizationID: "org-1",
					Status:         "Pending",
				},
			})
		case r.URL.Path == "/api/organization-invitations/inv-1" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/api/one-time-tokens" && r.Method == http.MethodPost:
			var req CreateOneTimeTokenRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			require.Equal(t, "test@example.com", req.Email)
			require.EqualValues(t, 3600, req.ExpiresIn)
			_ = json.NewEncoder(w).Encode(OneTimeTokenResponse{
				Token: "ott-mock-token-123",
			})
		case r.URL.Path == "/api/organizations/org-1/users" && r.Method == http.MethodGet:
			require.Equal(t, "1", r.URL.Query().Get("page"))
			require.Equal(t, "11", r.URL.Query().Get("page_size"))
			name := "Alice"
			email := "alice@example.com"
			_ = json.NewEncoder(w).Encode([]OrganizationUser{
				{
					ID:           "user-1",
					Name:         &name,
					PrimaryEmail: &email,
					CreatedAt:    1700000000,
				},
			})
		case r.URL.Path == "/api/organizations/org-1/users/user-2" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	conf := viper.New()
	conf.Set("security.logto.management.endpoint", server.URL)
	conf.Set("security.logto.management.app_id", "test-app")
	conf.Set("security.logto.management.app_secret", "test-secret")
	conf.Set("security.logto.management.resource", "https://default.logto.app/api")

	client, err := NewManagementClient(conf)
	require.NoError(t, err)

	// Test create invitation
	inv, err := client.CreateOrganizationInvitation(t.Context(), CreateInvitationRequest{
		Invitee:        "test@example.com",
		InviterID:      "user-1",
		OrganizationID: "org-1",
		ExpiresAt:      1700500000,
	})
	require.NoError(t, err)
	require.Equal(t, "inv-1", inv.ID)
	require.Equal(t, "test@example.com", inv.Invitee)

	// Test list invitations
	invList, err := client.ListOrganizationInvitations(t.Context(), "org-1")
	require.NoError(t, err)
	require.Len(t, invList, 1)
	require.Equal(t, "inv-1", invList[0].ID)

	// Test delete invitation
	err = client.DeleteOrganizationInvitation(t.Context(), "inv-1")
	require.NoError(t, err)

	// Test create one-time token
	ott, err := client.CreateOneTimeToken(t.Context(), CreateOneTimeTokenRequest{
		Email:     "test@example.com",
		ExpiresIn: 3600,
	})
	require.NoError(t, err)
	require.Equal(t, "ott-mock-token-123", ott.Token)

	// Test list members
	members, hasNext, err := client.ListOrganizationMembers(t.Context(), "org-1", "", 1, 10)
	require.NoError(t, err)
	require.Len(t, members, 1)
	require.False(t, hasNext)
	require.Equal(t, "user-1", members[0].ID)

	// Test remove member
	err = client.RemoveOrganizationMember(t.Context(), "org-1", "user-2")
	require.NoError(t, err)
}
