package repository

import (
	"aiot-backend/internal/model"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type OrganizationUserRepository interface {
	Create(ctx context.Context, orgUser *model.OrganizationUser) error
	ListByUser(ctx context.Context, userID string) ([]*model.OrganizationUser, error)
	ListOrgIDsByUser(ctx context.Context, userID string) ([]int64, error)
	IsMember(ctx context.Context, userID string, orgID int64) (bool, error)
	GetByUserAndOrg(ctx context.Context, userID string, orgID int64) (*model.OrganizationUser, error)
	UpdateLastLogin(ctx context.Context, userID string, orgID int64, loginTime time.Time) error
}

func NewOrganizationUserRepository(r *Repository) OrganizationUserRepository {
	return &organizationUserRepository{
		Repository: r,
	}
}

type organizationUserRepository struct {
	*Repository
}

func (r *organizationUserRepository) Create(ctx context.Context, orgUser *model.OrganizationUser) error {
	return r.DB(ctx).Create(orgUser).Error
}

func (r *organizationUserRepository) ListByUser(ctx context.Context, userID string) ([]*model.OrganizationUser, error) {
	var items []*model.OrganizationUser
	if err := r.DB(ctx).Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *organizationUserRepository) ListOrgIDsByUser(ctx context.Context, userID string) ([]int64, error) {
	var orgIDs []int64
	if err := r.DB(ctx).Model(&model.OrganizationUser{}).
		Where("user_id = ?", userID).
		Pluck("organization_id", &orgIDs).Error; err != nil {
		return nil, err
	}
	return orgIDs, nil
}

func (r *organizationUserRepository) IsMember(ctx context.Context, userID string, orgID int64) (bool, error) {
	var count int64
	if err := r.DB(ctx).Model(&model.OrganizationUser{}).
		Where("user_id = ? AND organization_id = ?", userID, orgID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *organizationUserRepository) GetByUserAndOrg(ctx context.Context, userID string, orgID int64) (*model.OrganizationUser, error) {
	var item model.OrganizationUser
	if err := r.DB(ctx).Where("user_id = ? AND organization_id = ?", userID, orgID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *organizationUserRepository) UpdateLastLogin(ctx context.Context, userID string, orgID int64, loginTime time.Time) error {
	return r.DB(ctx).Model(&model.OrganizationUser{}).
		Where("user_id = ? AND organization_id = ?", userID, orgID).
		Update("last_login_at", loginTime).Error
}
