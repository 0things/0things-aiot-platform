package repository

import (
	"context"
	"errors"
	"time"

	"0things-backend/internal/model"
	"gorm.io/gorm"
)

type RuleRepository struct {
	db *IoTDB
}

func NewRuleRepository(db *IoTDB) *RuleRepository {
	return &RuleRepository{db: db}
}

func (r *RuleRepository) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *RuleRepository) Find(ctx context.Context, id int64) (*model.Rule, error) {
	var rule model.Rule
	if err := r.DB(ctx).First(&rule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rule, nil
}

func (r *RuleRepository) List(ctx context.Context, page, size int, ruleType, status, search string) ([]model.Rule, int64, error) {
	query := r.DB(ctx).Model(&model.Rule{})
	if ruleType != "" {
		query = query.Where("type = ?", ruleType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rules []model.Rule
	if err := query.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&rules).Error; err != nil {
		return nil, 0, err
	}
	return rules, total, nil
}

func (r *RuleRepository) Create(ctx context.Context, rule *model.Rule) error {
	return r.DB(ctx).Create(rule).Error
}

func (r *RuleRepository) Save(ctx context.Context, rule *model.Rule) error {
	return r.DB(ctx).Save(rule).Error
}

func (r *RuleRepository) Delete(ctx context.Context, id int64) error {
	return r.DB(ctx).Delete(&model.Rule{}, id).Error
}

func (r *RuleRepository) UpdateStatus(ctx context.Context, rule *model.Rule, status string) error {
	return r.DB(ctx).Model(rule).Update("status", status).Error
}

func (r *RuleRepository) ListExecutions(ctx context.Context, ruleID int64, page, size int) ([]model.RuleExecution, int64, error) {
	query := r.DB(ctx).Model(&model.RuleExecution{}).Where("rule_id = ?", ruleID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var executions []model.RuleExecution
	if err := query.Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&executions).Error; err != nil {
		return nil, 0, err
	}
	return executions, total, nil
}

func (r *RuleRepository) CreateExecution(ctx context.Context, execution *model.RuleExecution) error {
	return r.DB(ctx).Create(execution).Error
}

func (r *RuleRepository) UpdateExecutionStats(ctx context.Context, rule *model.Rule, executedAt time.Time) error {
	return r.DB(ctx).Model(rule).Updates(map[string]any{
		"execution_count":       rule.ExecutionCount + 1,
		"success_count":         rule.SuccessCount + 1,
		"last_execution_status": "success",
		"last_executed_at":      executedAt,
	}).Error
}
