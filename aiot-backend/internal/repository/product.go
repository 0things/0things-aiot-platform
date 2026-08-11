package repository

import (
	"context"
	"errors"

	"0things-backend/internal/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *IoTDB
}

func NewProductRepository(db *IoTDB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *ProductRepository) Find(ctx context.Context, id int64) (*model.Product, error) {
	var product model.Product
	if err := r.DB(ctx).First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByKey(ctx context.Context, key string) (*model.Product, error) {
	var product model.Product
	if err := r.DB(ctx).Where("product_key = ?", key).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	return r.DB(ctx).Create(product).Error
}

func (r *ProductRepository) Save(ctx context.Context, product *model.Product) error {
	return r.DB(ctx).Save(product).Error
}

func (r *ProductRepository) CountDevices(ctx context.Context, productID int64) (int64, error) {
	var count int64
	err := r.DB(ctx).Model(&model.Device{}).Where("product_id = ?", productID).Count(&count).Error
	return count, err
}

func (r *ProductRepository) Delete(ctx context.Context, product *model.Product) error {
	return r.DB(ctx).Delete(product).Error
}

func (r *ProductRepository) List(ctx context.Context, page, size int, category, status, search string) ([]model.Product, int64, error) {
	query := r.DB(ctx).Model(&model.Product{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("product_key LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var products []model.Product
	if err := query.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductRepository) Restore(ctx context.Context, id int64) error {
	return r.DB(ctx).Unscoped().Model(&model.Product{}).Where("id = ?", id).Update("deleted_at", nil).Error
}
