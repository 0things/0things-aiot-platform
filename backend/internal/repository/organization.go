package repository

import (
	"aiot-backend/internal/model"
	"context"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *model.Organization) error
	GetByID(ctx context.Context, id int64) (*model.Organization, error)
	ListByIDs(ctx context.Context, ids []int64) ([]*model.Organization, error)
}

func NewOrganizationRepository(r *Repository) OrganizationRepository {
	return &organizationRepository{
		Repository: r,
	}
}

type organizationRepository struct {
	*Repository
}

func (r *organizationRepository) Create(ctx context.Context, org *model.Organization) error {
	return r.DB(ctx).Create(org).Error
}

func (r *organizationRepository) GetByID(ctx context.Context, id int64) (*model.Organization, error) {
	var org model.Organization
	if err := r.DB(ctx).Where("id = ?", id).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) ListByIDs(ctx context.Context, ids []int64) ([]*model.Organization, error) {
	if len(ids) == 0 {
		return []*model.Organization{}, nil
	}
	var orgs []*model.Organization
	if err := r.DB(ctx).Where("id IN ?", ids).Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}
