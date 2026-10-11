package logto

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type ManagementConfig struct {
	Endpoint          string
	AppID             string
	Secret            string
	Resource          string
	InvitationBaseURL string
}

type Organization struct {
	ID string `json:"id"`
}

type OrganizationInvitation struct {
	ID                  string   `json:"id"`
	InviterID           string   `json:"inviterId"`
	Invitee             string   `json:"invitee"`
	OrganizationID      string   `json:"organizationId"`
	Status              string   `json:"status"`
	CreatedAt           int64    `json:"createdAt"`
	ExpiresAt           int64    `json:"expiresAt"`
	OrganizationRoleIDs []string `json:"organizationRoleIds,omitempty"`
}

type CreateInvitationRequest struct {
	Invitee             string   `json:"invitee"`
	InviterID           string   `json:"inviterId"`
	OrganizationID      string   `json:"organizationId"`
	ExpiresAt           int64    `json:"expiresAt"`
	OrganizationRoleIDs []string `json:"organizationRoleIds,omitempty"`
	MessagePayload      any      `json:"messagePayload,omitempty"`
}

type OrganizationUser struct {
	ID           string  `json:"id"`
	Username     *string `json:"username,omitempty"`
	PrimaryEmail *string `json:"primaryEmail,omitempty"`
	PrimaryPhone *string `json:"primaryPhone,omitempty"`
	Name         *string `json:"name,omitempty"`
	Avatar       *string `json:"avatar,omitempty"`
	CreatedAt    int64   `json:"createdAt"`
}

type User struct {
	ID           string  `json:"id"`
	PrimaryEmail *string `json:"primaryEmail"`
}

type CreateOneTimeTokenRequest struct {
	Email     string `json:"email"`
	ExpiresIn int64  `json:"expiresIn,omitempty"`
}

type OneTimeTokenResponse struct {
	Token string `json:"token"`
}

type ManagementClient struct {
	config ManagementConfig
	client *http.Client

	mu          sync.Mutex
	accessToken string
	expiresAt   time.Time
}

func NewManagementClient(conf *viper.Viper) (*ManagementClient, error) {
	c := ManagementConfig{
		Endpoint:          strings.TrimRight(conf.GetString("security.logto.management.endpoint"), "/"),
		AppID:             strings.TrimSpace(conf.GetString("security.logto.management.app_id")),
		Secret:            strings.TrimSpace(conf.GetString("security.logto.management.app_secret")),
		Resource:          strings.TrimSpace(conf.GetString("security.logto.management.resource")),
		InvitationBaseURL: strings.TrimRight(strings.TrimSpace(conf.GetString("security.logto.invitation_base_url")), "/"),
	}
	if c.Endpoint == "" || c.AppID == "" || c.Secret == "" || c.Resource == "" {
		return nil, errors.New("security.logto.management endpoint, app_id, app_secret and resource are required")
	}
	return &ManagementClient{config: c, client: http.DefaultClient}, nil
}

func (c *ManagementClient) InvitationBaseURL() string {
	return c.config.InvitationBaseURL
}

func (c *ManagementClient) GetUser(ctx context.Context, userID string) (*User, error) {
	var user User
	if err := c.request(ctx, http.MethodGet, "/api/users/"+url.PathEscape(userID), nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *ManagementClient) UserOrganizations(ctx context.Context, userID string) ([]Organization, error) {
	var organizations []Organization
	if err := c.request(ctx, http.MethodGet, "/api/users/"+url.PathEscape(userID)+"/organizations", nil, &organizations); err != nil {
		return nil, err
	}
	return organizations, nil
}

func (c *ManagementClient) CreateOrganization(ctx context.Context, name string) (Organization, error) {
	var organization Organization
	body := map[string]string{"name": name}
	if err := c.request(ctx, http.MethodPost, "/api/organizations", body, &organization); err != nil {
		return Organization{}, err
	}
	return organization, nil
}

func (c *ManagementClient) AddUser(ctx context.Context, organizationID, userID string) error {
	return c.request(ctx, http.MethodPost, "/api/organizations/"+url.PathEscape(organizationID)+"/users", map[string][]string{"userIds": {userID}}, nil)
}

func (c *ManagementClient) CreateOrganizationInvitation(ctx context.Context, req CreateInvitationRequest) (*OrganizationInvitation, error) {
	var invitation OrganizationInvitation
	if err := c.request(ctx, http.MethodPost, "/api/organization-invitations", req, &invitation); err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (c *ManagementClient) CreateOneTimeToken(ctx context.Context, req CreateOneTimeTokenRequest) (*OneTimeTokenResponse, error) {
	var resp OneTimeTokenResponse
	if err := c.request(ctx, http.MethodPost, "/api/one-time-tokens", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *ManagementClient) ListOrganizationInvitations(ctx context.Context, organizationID string) ([]OrganizationInvitation, error) {
	var invitations []OrganizationInvitation
	path := "/api/organization-invitations?organizationId=" + url.QueryEscape(organizationID)
	if err := c.request(ctx, http.MethodGet, path, nil, &invitations); err != nil {
		return nil, err
	}
	return invitations, nil
}

func (c *ManagementClient) DeleteOrganizationInvitation(ctx context.Context, invitationID string) error {
	return c.request(ctx, http.MethodDelete, "/api/organization-invitations/"+url.PathEscape(invitationID), nil, nil)
}

func (c *ManagementClient) SendOrganizationInvitationMessage(ctx context.Context, invitationID, link string) error {
	body := map[string]string{"link": link}
	return c.request(ctx, http.MethodPost, "/api/organization-invitations/"+url.PathEscape(invitationID)+"/message", body, nil)
}

func (c *ManagementClient) GetOrganizationInvitation(ctx context.Context, invitationID string) (*OrganizationInvitation, error) {
	var invitation OrganizationInvitation
	if err := c.request(ctx, http.MethodGet, "/api/organization-invitations/"+url.PathEscape(invitationID), nil, &invitation); err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (c *ManagementClient) UpdateOrganizationInvitationStatus(ctx context.Context, invitationID, acceptedUserID string) error {
	body := map[string]string{"status": "Accepted", "acceptedUserId": acceptedUserID}
	return c.request(ctx, http.MethodPut, "/api/organization-invitations/"+url.PathEscape(invitationID)+"/status", body, nil)
}

func (c *ManagementClient) ListOrganizationMembers(ctx context.Context, organizationID, search string, page, pageSize int) ([]OrganizationUser, bool, error) {
	var users []OrganizationUser
	query := url.Values{}
	query.Set("page", fmt.Sprint(page))
	query.Set("page_size", fmt.Sprint(pageSize+1))
	if search != "" {
		query.Set("q", search)
	}
	path := "/api/organizations/" + url.PathEscape(organizationID) + "/users?" + query.Encode()
	if err := c.request(ctx, http.MethodGet, path, nil, &users); err != nil {
		return nil, false, err
	}
	hasNext := len(users) > pageSize
	if hasNext {
		users = users[:pageSize]
	}
	return users, hasNext, nil
}

func (c *ManagementClient) RemoveOrganizationMember(ctx context.Context, organizationID, userID string) error {
	return c.request(ctx, http.MethodDelete, "/api/organizations/"+url.PathEscape(organizationID)+"/users/"+url.PathEscape(userID), nil, nil)
}

func (c *ManagementClient) request(ctx context.Context, method, path string, body any, result any) error {
	token, err := c.token(ctx)
	if err != nil {
		return err
	}
	var requestBody *strings.Reader
	if body == nil {
		requestBody = strings.NewReader("")
	} else {
		payload, marshalErr := json.Marshal(body)
		if marshalErr != nil {
			return fmt.Errorf("marshal Logto request: %w", marshalErr)
		}
		requestBody = strings.NewReader(string(payload))
	}
	req, err := http.NewRequestWithContext(ctx, method, c.config.Endpoint+path, requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("call Logto Management API: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("Logto Management API returned %s", resp.Status)
	}
	if result != nil && resp.StatusCode != http.StatusNoContent && resp.ContentLength != 0 {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode Logto Management API response: %w", err)
		}
	}
	return nil
}

func (c *ManagementClient) token(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.accessToken != "" && time.Until(c.expiresAt) > time.Minute {
		return c.accessToken, nil
	}
	form := url.Values{
		"grant_type": {"client_credentials"},
		"resource":   {c.config.Resource},
		"scope":      {"all"},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.Endpoint+"/oidc/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(c.config.AppID, c.config.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request Logto Management token: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Logto Management token returned %s", resp.Status)
	}
	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode Logto Management token: %w", err)
	}
	if result.AccessToken == "" {
		return "", errors.New("Logto Management token is empty")
	}
	c.accessToken = result.AccessToken
	c.expiresAt = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	return c.accessToken, nil
}
