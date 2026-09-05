package repository

import (
	"context"
	"errors"

	"aiot-backend/internal/model"

	"gorm.io/gorm"
)

type ProductTSLRepository struct {
	db *gorm.DB
}

func NewProductTSLRepository(db *gorm.DB) *ProductTSLRepository {
	return &ProductTSLRepository{db: db}
}

func (r *ProductTSLRepository) FindByProductID(ctx context.Context, productID int64) (*model.ProductTSL, error) {
	var tsl model.ProductTSL
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).First(&tsl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &tsl, nil
}

func (r *ProductTSLRepository) Save(ctx context.Context, tsl *model.ProductTSL) error {
	return r.db.WithContext(ctx).Save(tsl).Error
}

func (r *ProductTSLRepository) Create(ctx context.Context, tsl *model.ProductTSL) error {
	return r.db.WithContext(ctx).Create(tsl).Error
}

func (r *ProductTSLRepository) Delete(ctx context.Context, tsl *model.ProductTSL) error {
	return r.db.WithContext(ctx).Delete(tsl).Error
}
