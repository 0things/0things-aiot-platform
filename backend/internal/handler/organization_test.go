package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"aiot-backend/internal/middleware"
	"aiot-backend/internal/service"
	"aiot-backend/internal/tenant"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type mockOrgService struct {
	createOrgFn func(ctx context.Context, userID, name string) (string, error)
	createInvFn func(ctx context.Context, organizationID, inviterID, inviteeEmail string) (*service.OrganizationInvitation, error)
	getInvFn    func(ctx context.Context, invitationID, userID string) (*service.OrganizationInvitation, error)
	acceptInvFn func(ctx context.Context, invitationID, userID string) error
	listInvsFn  func(ctx context.Context, organizationID string, page, pageSize int) ([]service.OrganizationInvitation, int, int, error)
	deleteInvFn func(ctx context.Context, organizationID, invitationID string) error
	listMemsFn  func(ctx context.Context, organizationID, query string, page, pageSize int) ([]service.OrganizationMember, bool, error)
	removeMemFn func(ctx context.Context, organizationID, targetUserID string) error
}

func (m *mockOrgService) CreateOrganization(ctx context.Context, userID, name string) (string, error) {
	if m.createOrgFn != nil {
		return m.createOrgFn(ctx, userID, name)
	}
	return "org-1", nil
}

func (m *mockOrgService) CreateInvitation(ctx context.Context, organizationID, inviterID, inviteeEmail string) (*service.OrganizationInvitation, error) {
	if m.createInvFn != nil {
		return m.createInvFn(ctx, organizationID, inviterID, inviteeEmail)
	}
	return &service.OrganizationInvitation{
		ID:             "inv-1",
		InviterID:      inviterID,
		Invitee:        inviteeEmail,
		OrganizationID: organizationID,
		Status:         "Pending",
	}, nil
}

func (m *mockOrgService) GetInvitation(ctx context.Context, invitationID, userID string) (*service.OrganizationInvitation, error) {
	if m.getInvFn != nil {
		return m.getInvFn(ctx, invitationID, userID)
	}
	return &service.OrganizationInvitation{
		ID:             invitationID,
		Invitee:        "test@example.com",
		OrganizationID: "org-1",
		Status:         "Pending",
	}, nil
}

func (m *mockOrgService) AcceptInvitation(ctx context.Context, invitationID, userID string) error {
	if m.acceptInvFn != nil {
		return m.acceptInvFn(ctx, invitationID, userID)
	}
	return nil
}

func (m *mockOrgService) ListInvitations(ctx context.Context, organizationID string, page, pageSize int) ([]service.OrganizationInvitation, int, int, error) {
	if m.listInvsFn != nil {
		return m.listInvsFn(ctx, organizationID, page, pageSize)
	}
	return []service.OrganizationInvitation{
		{ID: "inv-1", Invitee: "test@example.com", Status: "Pending"},
	}, 1, 1, nil
}

func (m *mockOrgService) DeleteInvitation(ctx context.Context, organizationID, invitationID string) error {
	if m.deleteInvFn != nil {
		return m.deleteInvFn(ctx, organizationID, invitationID)
	}
	return nil
}

func (m *mockOrgService) ListMembers(ctx context.Context, organizationID, query string, page, pageSize int) ([]service.OrganizationMember, bool, error) {
	if m.listMemsFn != nil {
		return m.listMemsFn(ctx, organizationID, query, page, pageSize)
	}
	name := "Alice"
	email := "alice@example.com"
	return []service.OrganizationMember{
		{ID: "user-1", Name: &name, PrimaryEmail: &email},
	}, false, nil
}

func (m *mockOrgService) RemoveMember(ctx context.Context, organizationID, targetUserID string) error {
	if m.removeMemFn != nil {
		return m.removeMemFn(ctx, organizationID, targetUserID)
	}
	return nil
}

func setupOrgTestRouter(svc *mockOrgService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewOrganizationHandler(&Handler{}, svc)

	// middleware to mock auth context
	r.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, "user-1")
		c.Set(string(tenant.OrganizationKey), "org-1")
		c.Next()
	})

	r.POST("/organizations", h.CreateOrganization)
	r.POST("/organizations/invitations", h.CreateInvitation)
	r.GET("/organizations/invitations", h.ListInvitations)
	r.GET("/organizations/invitations/:id", h.GetInvitation)
	r.POST("/organizations/invitations/:id/accept", h.AcceptInvitation)
	r.DELETE("/organizations/invitations/:id", h.DeleteInvitation)
	r.GET("/organizations/members", h.ListMembers)
	r.DELETE("/organizations/members/:userId", h.RemoveMember)

	return r
}

func TestOrganizationHandler_Endpoints(t *testing.T) {
	svc := &mockOrgService{}
	r := setupOrgTestRouter(svc)

	// 1. Test explicit organization creation
	svc.createOrgFn = func(_ context.Context, userID, name string) (string, error) {
		require.Equal(t, "user-1", userID)
		require.Equal(t, "Acme", name)
		return "org-1", nil
	}
	req := httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBufferString(`{"name":"Acme"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"organizationId":"org-1"`)

	// 2. Test CreateInvitation
	body, _ := json.Marshal(map[string]string{"invitee": "bob@example.com"})
	req = httptest.NewRequest(http.MethodPost, "/organizations/invitations", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"id":"inv-1"`)
	require.Contains(t, w.Body.String(), `"invitee":"bob@example.com"`)

	// 3. Test ListInvitations
	req = httptest.NewRequest(http.MethodGet, "/organizations/invitations", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"inv-1"`)

	// 4. Test DeleteInvitation
	req = httptest.NewRequest(http.MethodDelete, "/organizations/invitations/inv-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// 5. Test ListMembers
	req = httptest.NewRequest(http.MethodGet, "/organizations/members", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"alice@example.com"`)

	// 6. Test RemoveMember
	req = httptest.NewRequest(http.MethodDelete, "/organizations/members/user-2", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// 7. Test GetInvitation
	req = httptest.NewRequest(http.MethodGet, "/organizations/invitations/inv-1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"inv-1"`)

	// 8. Test AcceptInvitation
	req = httptest.NewRequest(http.MethodPost, "/organizations/invitations/inv-1/accept", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCreateOrganizationRequiresName(t *testing.T) {
	r := setupOrgTestRouter(&mockOrgService{})
	for _, body := range []string{`{"name":""}`, `{"name":"   "}`} {
		req := httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	}
}

func TestOrganizationHandlerRejectsWrongInvitee(t *testing.T) {
	svc := &mockOrgService{
		getInvFn: func(context.Context, string, string) (*service.OrganizationInvitation, error) {
			return nil, service.ErrInvitationForbidden
		},
		acceptInvFn: func(context.Context, string, string) error {
			return service.ErrInvitationForbidden
		},
	}
	r := setupOrgTestRouter(svc)
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		req := httptest.NewRequest(method, "/organizations/invitations/inv-1", nil)
		if method == http.MethodPost {
			req.URL.Path += "/accept"
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusForbidden, w.Code)
	}
}

func TestInvitationResponseFormatsTimesInUTCPlusEight(t *testing.T) {
	createdAt := time.Date(2026, time.October, 10, 2, 50, 28, 0, time.UTC).UnixMilli()
	expiresAt := time.Date(2026, time.October, 17, 2, 50, 27, 0, time.UTC).UnixMilli()
	resp := toInvitationResp(&service.OrganizationInvitation{
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	})

	require.Equal(t, "2026-10-10 10:50:28", resp.CreatedAt)
	require.Equal(t, "2026-10-17 10:50:27", resp.ExpiresAt)
	require.Empty(t, formatInvitationTime(0))
}

func TestInvitationResponseMarksExpired(t *testing.T) {
	resp := toInvitationResp(&service.OrganizationInvitation{
		ExpiresAt: time.Now().Add(-time.Minute).UnixMilli(),
	})
	require.True(t, resp.IsExpired)
}
