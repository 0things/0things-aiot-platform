package repository

import (
	"aiot-backend/internal/model"
	"context"
)

type CategoryRepository struct{ db *IoTDB }

func NewCategoryRepository(db *IoTDB) *CategoryRepository { return &CategoryRepository{db: db} }

func (r *CategoryRepository) List(ctx context.Context) ([]model.Category, error) {
	var items []model.Category
	err := r.db.WithContext(ctx).Where("enabled = ?", true).Order("sort, id").Find(&items).Error
	return items, err
}
