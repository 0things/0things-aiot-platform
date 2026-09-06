package repository

import (
	"context"

	"aiot-backend/internal/model"

	"gorm.io/gorm"
)

type RuleNodeDefinitionRepository struct {
	db *gorm.DB
}

func NewRuleNodeDefinitionRepository(db *gorm.DB) *RuleNodeDefinitionRepository {
	return &RuleNodeDefinitionRepository{db: db}
}

func (r *RuleNodeDefinitionRepository) ListEnabled(ctx context.Context) ([]model.RuleNodeDefinition, error) {
	var definitions []model.RuleNodeDefinition
	err := r.db.WithContext(ctx).
		Where("enabled = ?", true).
		Order("is_system DESC, category ASC, sort_order ASC, name ASC").
		Find(&definitions).Error
	return definitions, err
}
