package repository

import (
	"context"
	"errors"

	"aiot-backend/internal/model"
	"aiot-backend/internal/tenant"

	"gorm.io/gorm"
)

type RuleChainRepository struct{ db *gorm.DB }

func NewRuleChainRepository(db *gorm.DB) *RuleChainRepository { return &RuleChainRepository{db: db} }

func (r *RuleChainRepository) Find(ctx context.Context, uuid string) (*model.RuleChain, error) {
	var chain model.RuleChain
	err := r.db.WithContext(ctx).Where("uuid = ? AND organization_id = ?", uuid, tenant.GetOrganizationID(ctx)).First(&chain).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &chain, nil
}

func (r *RuleChainRepository) List(ctx context.Context, page, size int, search, status string) ([]model.RuleChain, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.RuleChain{}).Where("organization_id = ?", tenant.GetOrganizationID(ctx))
	if search != "" {
		q = q.Where("name LIKE ?", "%"+search+"%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.RuleChain
	if err := q.Order("updated_at DESC").Offset((page - 1) * size).Limit(size).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *RuleChainRepository) Create(ctx context.Context, chain *model.RuleChain) error {
	chain.OrganizationID = tenant.GetOrganizationID(ctx)
	return r.db.WithContext(ctx).Create(chain).Error
}

func (r *RuleChainRepository) NameExists(ctx context.Context, name, excludeUUID string) (bool, error) {
	q := r.db.WithContext(ctx).Model(&model.RuleChain{}).
		Where("organization_id = ? AND name = ?", tenant.GetOrganizationID(ctx), name)
	if excludeUUID != "" {
		q = q.Where("uuid <> ?", excludeUUID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RuleChainRepository) Save(ctx context.Context, chain *model.RuleChain) error {
	if _, err := r.Find(ctx, chain.UUID); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(chain).Error
}

func (r *RuleChainRepository) Delete(ctx context.Context, uuid string) error {
	chain, err := r.Find(ctx, uuid)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(chain).Error
}
