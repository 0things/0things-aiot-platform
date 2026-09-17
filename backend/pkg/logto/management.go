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
	Endpoint string
	AppID    string
	Secret   string
	Resource string
}

type Organization struct {
	ID string `json:"id"`
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
		Endpoint: strings.TrimRight(conf.GetString("security.logto.management.endpoint"), "/"),
		AppID:    strings.TrimSpace(conf.GetString("security.logto.management.app_id")),
		Secret:   strings.TrimSpace(conf.GetString("security.logto.management.app_secret")),
		Resource: strings.TrimSpace(conf.GetString("security.logto.management.resource")),
	}
	if c.Endpoint == "" || c.AppID == "" || c.Secret == "" || c.Resource == "" {
		return nil, errors.New("security.logto.management endpoint, app_id, app_secret and resource are required")
	}
	return &ManagementClient{config: c, client: http.DefaultClient}, nil
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
	if result != nil && resp.ContentLength != 0 {
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
