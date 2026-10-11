package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"aiot-backend/pkg/logto"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestOrganizationProvisioningService_InvitationsAndMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/oidc/token":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token": "mock-token",
				"expires_in":   3600,
			})
		case r.URL.Path == "/api/organization-invitations" && r.Method == http.MethodPost:
			var req logto.CreateInvitationRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			require.Empty(t, req.OrganizationRoleIDs)
			_ = json.NewEncoder(w).Encode(logto.OrganizationInvitation{
				ID:             "inv-123",
				InviterID:      req.InviterID,
				Invitee:        req.Invitee,
				OrganizationID: req.OrganizationID,
				Status:         "Pending",
				CreatedAt:      1700000000,
				ExpiresAt:      req.ExpiresAt,
			})
		case r.URL.Path == "/api/organization-invitations" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode([]logto.OrganizationInvitation{
				{
					ID:             "inv-123",
					InviterID:      "user-1",
					Invitee:        "bob@example.com",
					OrganizationID: "org-1",
					Status:         "Pending",
				},
			})
		case r.URL.Path == "/api/one-time-tokens" && r.Method == http.MethodPost:
			var req logto.CreateOneTimeTokenRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			require.Equal(t, "bob@example.com", req.Email)
			require.EqualValues(t, 60*60, req.ExpiresIn)
			_ = json.NewEncoder(w).Encode(logto.OneTimeTokenResponse{
				Token: "ott-mock-123",
			})
		case r.URL.Path == "/api/organization-invitations/inv-123/message" && r.Method == http.MethodPost:
			var payload map[string]string
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "http://localhost:5173/invitations/accept?invitationId=inv-123&token=ott-mock-123&email=bob%40example.com", payload["link"])
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/api/users/user-2" && r.Method == http.MethodGet:
			email := "bob@example.com"
			_ = json.NewEncoder(w).Encode(logto.User{ID: "user-2", PrimaryEmail: &email})
		case r.URL.Path == "/api/users/user-other" && r.Method == http.MethodGet:
			email := "other@example.com"
			_ = json.NewEncoder(w).Encode(logto.User{ID: "user-other", PrimaryEmail: &email})
		case r.URL.Path == "/api/organization-invitations/inv-123" && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode(logto.OrganizationInvitation{
				ID:             "inv-123",
				InviterID:      "user-1",
				Invitee:        "bob@example.com",
				OrganizationID: "org-1",
				Status:         "Pending",
				ExpiresAt:      time.Now().Add(time.Hour).UnixMilli(),
			})
		case r.URL.Path == "/api/organization-invitations/inv-123/status" && r.Method == http.MethodPut:
			var payload map[string]string
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "Accepted", payload["status"])
			require.Equal(t, "user-2", payload["acceptedUserId"])
			w.WriteHeader(http.StatusOK)
		case r.URL.Path == "/api/organization-invitations/inv-123" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/api/organizations/org-1/users" && r.Method == http.MethodGet:
			name := "Bob"
			email := "bob@example.com"
			_ = json.NewEncoder(w).Encode([]logto.OrganizationUser{
				{
					ID:           "user-2",
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
	conf.Set("security.logto.invitation_base_url", "http://localhost:5173")

	client, err := logto.NewManagementClient(conf)
	require.NoError(t, err)

	svc := NewOrganizationProvisioningService(client)

	// Test validation
	_, err = svc.CreateInvitation(t.Context(), "", "user-1", "bob@example.com")
	require.Error(t, err)
	_, err = svc.CreateInvitation(t.Context(), "org-1", "user-1", "")
	require.Error(t, err)

	// Test CreateInvitation
	inv, err := svc.CreateInvitation(t.Context(), "org-1", "user-1", "bob@example.com")
	require.NoError(t, err)
	require.Equal(t, "inv-123", inv.ID)
	require.Equal(t, "bob@example.com", inv.Invitee)

	// Test GetInvitation
	gotInv, err := svc.GetInvitation(t.Context(), "inv-123", "user-2")
	require.NoError(t, err)
	require.Equal(t, "inv-123", gotInv.ID)
	_, err = svc.GetInvitation(t.Context(), "inv-123", "user-other")
	require.ErrorIs(t, err, ErrInvitationForbidden)
	err = svc.AcceptInvitation(t.Context(), "inv-123", "user-other")
	require.ErrorIs(t, err, ErrInvitationForbidden)

	// Test AcceptInvitation
	err = svc.AcceptInvitation(t.Context(), "inv-123", "user-2")
	require.NoError(t, err)

	// Test ListInvitations
	invs, total, pendingCount, err := svc.ListInvitations(t.Context(), "org-1", 1, 10)
	require.NoError(t, err)
	require.Len(t, invs, 1)
	require.Equal(t, 1, total)
	require.Equal(t, 1, pendingCount)

	// Test DeleteInvitation
	err = svc.DeleteInvitation(t.Context(), "org-1", "inv-123")
	require.NoError(t, err)

	// Test ListMembers
	members, hasNext, err := svc.ListMembers(t.Context(), "org-1", "", 1, 10)
	require.NoError(t, err)
	require.Len(t, members, 1)
	require.False(t, hasNext)
	require.Equal(t, "user-2", members[0].ID)

	// Test RemoveMember
	err = svc.RemoveMember(t.Context(), "org-1", "user-2")
	require.NoError(t, err)
}

func TestCreateOrganizationUsesRequestedNameAndAddsCurrentUser(t *testing.T) {
	created := false
	added := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oidc/token":
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "mock-token", "expires_in": 3600})
		case r.URL.Path == "/api/organizations" && r.Method == http.MethodPost:
			var payload map[string]string
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "Acme", payload["name"])
			created = true
			_ = json.NewEncoder(w).Encode(logto.Organization{ID: "org-new"})
		case r.URL.Path == "/api/organizations/org-new/users" && r.Method == http.MethodPost:
			var payload map[string][]string
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, []string{"user-1"}, payload["userIds"])
			added = true
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
	client, err := logto.NewManagementClient(conf)
	require.NoError(t, err)
	svc := NewOrganizationProvisioningService(client)
	_, err = svc.CreateOrganization(t.Context(), "user-1", " ")
	require.Error(t, err)
	require.False(t, created)
	id, err := svc.CreateOrganization(t.Context(), "user-1", " Acme ")
	require.NoError(t, err)
	require.Equal(t, "org-new", id)
	require.True(t, created)
	require.True(t, added)
}

func TestCreateInvitationReportsMailFailure(t *testing.T) {
	deleted := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oidc/token":
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "mock-token", "expires_in": 3600})
		case r.URL.Path == "/api/organization-invitations" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(logto.OrganizationInvitation{ID: "inv-123"})
		case r.URL.Path == "/api/one-time-tokens" && r.Method == http.MethodPost:
			_ = json.NewEncoder(w).Encode(logto.OneTimeTokenResponse{Token: "ott-mock-123"})
		case r.URL.Path == "/api/organization-invitations/inv-123/message":
			http.Error(w, "mail failed", http.StatusBadGateway)
		case r.URL.Path == "/api/organization-invitations/inv-123" && r.Method == http.MethodDelete:
			deleted = true
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
	conf.Set("security.logto.invitation_base_url", "http://localhost:5173")
	client, err := logto.NewManagementClient(conf)
	require.NoError(t, err)
	_, err = NewOrganizationProvisioningService(client).CreateInvitation(t.Context(), "org-1", "user-1", "bob@example.com")
	require.Error(t, err)
	require.True(t, deleted)
}

func TestCreateInvitationReportsOneTimeTokenFailure(t *testing.T) {
	deleted := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/oidc/token":
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "mock-token", "expires_in": 3600})
		case r.URL.Path == "/api/organization-invitations" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(logto.OrganizationInvitation{ID: "inv-123"})
		case r.URL.Path == "/api/one-time-tokens" && r.Method == http.MethodPost:
			http.Error(w, "token failed", http.StatusBadGateway)
		case r.URL.Path == "/api/organization-invitations/inv-123" && r.Method == http.MethodDelete:
			deleted = true
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
	conf.Set("security.logto.invitation_base_url", "http://localhost:5173")
	client, err := logto.NewManagementClient(conf)
	require.NoError(t, err)
	_, err = NewOrganizationProvisioningService(client).CreateInvitation(t.Context(), "org-1", "user-1", "bob@example.com")
	require.ErrorContains(t, err, "create invitation one-time token")
	require.True(t, deleted)
}
