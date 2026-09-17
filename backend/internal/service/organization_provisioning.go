package service

import (
	"context"
	"fmt"

	"aiot-backend/pkg/logto"

	"golang.org/x/sync/singleflight"
)

type OrganizationProvisioningServiceInterface interface {
	EnsureOrganization(ctx context.Context, userID string) (string, error)
}

type OrganizationProvisioningService struct {
	client *logto.ManagementClient
	sf     singleflight.Group
}

func NewOrganizationProvisioningService(client *logto.ManagementClient) *OrganizationProvisioningService {
	return &OrganizationProvisioningService{client: client}
}

func (s *OrganizationProvisioningService) EnsureOrganization(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("user id is required")
	}

	res, err, _ := s.sf.Do(userID, func() (interface{}, error) {
		organizations, err := s.client.UserOrganizations(ctx, userID)
		if err != nil {
			return "", err
		}
		if len(organizations) > 0 {
			return organizations[0].ID, nil
		}
		organization, err := s.client.CreateOrganization(ctx, "Personal organization "+userID)
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
	})
	if err != nil {
		return "", err
	}
	return res.(string), nil
}
