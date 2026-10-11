package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"aiot-backend/pkg/logto"
)

type OrganizationProvisioningServiceInterface interface {
	CreateOrganization(ctx context.Context, userID, name string) (string, error)
	CreateInvitation(ctx context.Context, organizationID, inviterID, inviteeEmail string) (*OrganizationInvitation, error)
	GetInvitation(ctx context.Context, invitationID, userID string) (*OrganizationInvitation, error)
	AcceptInvitation(ctx context.Context, invitationID, userID string) error
	ListInvitations(ctx context.Context, organizationID string, page, pageSize int) ([]OrganizationInvitation, int, int, error)
	DeleteInvitation(ctx context.Context, organizationID, invitationID string) error
	ListMembers(ctx context.Context, organizationID, query string, page, pageSize int) ([]OrganizationMember, bool, error)
	RemoveMember(ctx context.Context, organizationID, targetUserID string) error
}

type OrganizationInvitation struct {
	ID                  string
	InviterID           string
	Invitee             string
	OrganizationID      string
	Status              string
	CreatedAt           int64
	ExpiresAt           int64
	OrganizationRoleIDs []string
}

type OrganizationMember struct {
	ID           string
	Username     *string
	PrimaryEmail *string
	PrimaryPhone *string
	Name         *string
	Avatar       *string
	CreatedAt    int64
}

var ErrInvitationForbidden = errors.New("invitation does not belong to the signed-in email")
var ErrInvitationUnavailable = errors.New("invitation is no longer available")

type OrganizationProvisioningService struct {
	client *logto.ManagementClient
}

func NewOrganizationProvisioningService(client *logto.ManagementClient) *OrganizationProvisioningService {
	return &OrganizationProvisioningService{client: client}
}

func (s *OrganizationProvisioningService) CreateOrganization(ctx context.Context, userID, name string) (string, error) {
	name = strings.TrimSpace(name)
	if userID == "" || name == "" {
		return "", fmt.Errorf("user id and organization name are required")
	}
	organization, err := s.client.CreateOrganization(ctx, name)
	if err != nil {
		return "", err
	}
	if organization.ID == "" {
		return "", fmt.Errorf("Logto created an organization without an id")
	}
	if err := s.client.AddUser(ctx, organization.ID, userID); err != nil {
		return "", err
	}
	return organization.ID, nil
}

func (s *OrganizationProvisioningService) CreateInvitation(ctx context.Context, organizationID, inviterID, inviteeEmail string) (*OrganizationInvitation, error) {
	if organizationID == "" || inviterID == "" || inviteeEmail == "" {
		return nil, fmt.Errorf("organization id, inviter id and invitee email are required")
	}
	baseURL := s.client.InvitationBaseURL()
	if baseURL == "" {
		return nil, fmt.Errorf("security.logto.invitation_base_url is required")
	}
	invitationValidity := 7 * 24 * time.Hour
	// Logto expects invitation expiration as epoch milliseconds.
	expiresAt := time.Now().Add(invitationValidity).UnixMilli()
	req := logto.CreateInvitationRequest{
		Invitee:        inviteeEmail,
		InviterID:      inviterID,
		OrganizationID: organizationID,
		ExpiresAt:      expiresAt,
		MessagePayload: false,
	}
	inv, err := s.client.CreateOrganizationInvitation(ctx, req)
	if err != nil {
		return nil, err
	}

	ott, err := s.client.CreateOneTimeToken(ctx, logto.CreateOneTimeTokenRequest{
		Email:     inviteeEmail,
		ExpiresIn: int64(time.Hour.Seconds()),
	})
	if err != nil {
		_ = s.client.DeleteOrganizationInvitation(ctx, inv.ID)
		return nil, fmt.Errorf("create invitation one-time token: %w", err)
	}
	if ott.Token == "" {
		_ = s.client.DeleteOrganizationInvitation(ctx, inv.ID)
		return nil, fmt.Errorf("Logto created an empty one-time token")
	}

	link := baseURL + "/invitations/accept?invitationId=" + url.QueryEscape(inv.ID) +
		"&token=" + url.QueryEscape(ott.Token) + "&email=" + url.QueryEscape(inviteeEmail)
	if err := s.client.SendOrganizationInvitationMessage(ctx, inv.ID, link); err != nil {
		_ = s.client.DeleteOrganizationInvitation(ctx, inv.ID)
		return nil, fmt.Errorf("send organization invitation: %w", err)
	}
	return mapOrganizationInvitation(inv), nil
}

func (s *OrganizationProvisioningService) GetInvitation(ctx context.Context, invitationID, userID string) (*OrganizationInvitation, error) {
	if invitationID == "" || userID == "" {
		return nil, fmt.Errorf("invitation id and user id are required")
	}
	inv, err := s.client.GetOrganizationInvitation(ctx, invitationID)
	if err != nil {
		return nil, err
	}
	user, err := s.client.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.PrimaryEmail == nil || !strings.EqualFold(*user.PrimaryEmail, inv.Invitee) {
		return nil, ErrInvitationForbidden
	}
	return mapOrganizationInvitation(inv), nil
}

func (s *OrganizationProvisioningService) AcceptInvitation(ctx context.Context, invitationID, userID string) error {
	inv, err := s.GetInvitation(ctx, invitationID, userID)
	if err != nil {
		return err
	}
	if inv.Status != "Pending" || inv.ExpiresAt <= time.Now().UnixMilli() {
		return ErrInvitationUnavailable
	}
	return s.client.UpdateOrganizationInvitationStatus(ctx, invitationID, userID)
}

func (s *OrganizationProvisioningService) ListInvitations(ctx context.Context, organizationID string, page, pageSize int) ([]OrganizationInvitation, int, int, error) {
	if organizationID == "" {
		return nil, 0, 0, fmt.Errorf("organization id is required")
	}
	invitations, err := s.client.ListOrganizationInvitations(ctx, organizationID)
	if err != nil {
		return nil, 0, 0, err
	}
	total := len(invitations)
	pendingCount := 0
	for _, invitation := range invitations {
		if strings.EqualFold(invitation.Status, "pending") {
			pendingCount++
		}
	}
	start := (page - 1) * pageSize
	if start >= total {
		return []OrganizationInvitation{}, total, pendingCount, nil
	}
	end := min(start+pageSize, total)
	items := make([]OrganizationInvitation, 0, end-start)
	for _, invitation := range invitations[start:end] {
		items = append(items, *mapOrganizationInvitation(&invitation))
	}
	return items, total, pendingCount, nil
}

func (s *OrganizationProvisioningService) DeleteInvitation(ctx context.Context, organizationID, invitationID string) error {
	if organizationID == "" || invitationID == "" {
		return fmt.Errorf("organization id and invitation id are required")
	}
	inv, err := s.client.GetOrganizationInvitation(ctx, invitationID)
	if err != nil {
		return err
	}
	if inv.OrganizationID != organizationID {
		return ErrInvitationForbidden
	}
	return s.client.DeleteOrganizationInvitation(ctx, invitationID)
}

func (s *OrganizationProvisioningService) ListMembers(ctx context.Context, organizationID, query string, page, pageSize int) ([]OrganizationMember, bool, error) {
	if organizationID == "" {
		return nil, false, fmt.Errorf("organization id is required")
	}
	users, hasNext, err := s.client.ListOrganizationMembers(ctx, organizationID, query, page, pageSize)
	if err != nil {
		return nil, false, err
	}
	members := make([]OrganizationMember, len(users))
	for i, user := range users {
		members[i] = OrganizationMember{
			ID: user.ID, Username: user.Username, PrimaryEmail: user.PrimaryEmail,
			PrimaryPhone: user.PrimaryPhone, Name: user.Name, Avatar: user.Avatar,
			CreatedAt: user.CreatedAt,
		}
	}
	return members, hasNext, nil
}

func mapOrganizationInvitation(invitation *logto.OrganizationInvitation) *OrganizationInvitation {
	if invitation == nil {
		return nil
	}
	return &OrganizationInvitation{
		ID: invitation.ID, InviterID: invitation.InviterID, Invitee: invitation.Invitee,
		OrganizationID: invitation.OrganizationID, Status: invitation.Status,
		CreatedAt: invitation.CreatedAt, ExpiresAt: invitation.ExpiresAt,
		OrganizationRoleIDs: invitation.OrganizationRoleIDs,
	}
}

func (s *OrganizationProvisioningService) RemoveMember(ctx context.Context, organizationID, targetUserID string) error {
	if organizationID == "" || targetUserID == "" {
		return fmt.Errorf("organization id and target user id are required")
	}
	return s.client.RemoveOrganizationMember(ctx, organizationID, targetUserID)
}
