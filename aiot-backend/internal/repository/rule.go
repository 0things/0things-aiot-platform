package repository

import (
	"context"
	"errors"
	"time"

	"0things-backend/internal/model"
	"0things-backend/internal/tenant"
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
	q := useIoTQuery(r.db)
	rule, err := q.Rule.WithContext(ctx).Join(q.Product, q.Product.ID.EqCol(q.Rule.ProductID)).Where(q.Rule.ID.Eq(id), q.Product.TenantID.Eq(tenant.GetTenantID(ctx))).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return rule, nil
}

func (r *RuleRepository) List(ctx context.Context, page, size int, ruleType, status, search string) ([]model.Rule, int64, error) {
	q := useIoTQuery(r.db)
	rules := q.Rule.WithContext(ctx).Join(q.Product, q.Product.ID.EqCol(q.Rule.ProductID)).Where(q.Product.TenantID.Eq(tenant.GetTenantID(ctx)))
	if ruleType != "" {
		rules = rules.Where(q.Rule.Type.Eq(ruleType))
	}
	if status != "" {
		rules = rules.Where(q.Rule.Status.Eq(status))
	}
	if search != "" {
		rules = rules.Where(q.Rule.Name.Like("%" + search + "%"))
	}
	items, total, err := rules.Order(q.Rule.CreatedAt.Desc()).FindByPage((page-1)*size, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]model.Rule, len(items))
	for i := range items {
		result[i] = *items[i]
	}
	return result, total, nil
}

func (r *RuleRepository) Create(ctx context.Context, rule *model.Rule) error {
	q := useIoTQuery(r.db)
	if _, err := q.Product.WithContext(ctx).Where(q.Product.ID.Eq(rule.ProductID), q.Product.TenantID.Eq(tenant.GetTenantID(ctx))).First(); err != nil {
		return err
	}
	return q.Rule.WithContext(ctx).Create(rule)
}

func (r *RuleRepository) Save(ctx context.Context, rule *model.Rule) error {
	if _, err := r.Find(ctx, rule.ID); err != nil {
		return err
	}
	return useIoTQuery(r.db).Rule.WithContext(ctx).Save(rule)
}

func (r *RuleRepository) Delete(ctx context.Context, id int64) error {
	rule, err := r.Find(ctx, id)
	if err != nil {
		return err
	}
	_, err = useIoTQuery(r.db).Rule.WithContext(ctx).Delete(rule)
	return err
}

func (r *RuleRepository) UpdateStatus(ctx context.Context, rule *model.Rule, status string) error {
	_, err := useIoTQuery(r.db).Rule.WithContext(ctx).Where(useIoTQuery(r.db).Rule.ID.Eq(rule.ID)).UpdateSimple(useIoTQuery(r.db).Rule.Status.Value(status))
	return err
}

func (r *RuleRepository) ListExecutions(ctx context.Context, ruleID int64, page, size int) ([]model.RuleExecution, int64, error) {
	q := useIoTQuery(r.db)
	executions := q.RuleExecution.WithContext(ctx).Join(q.Rule, q.Rule.ID.EqCol(q.RuleExecution.RuleID)).Join(q.Product, q.Product.ID.EqCol(q.Rule.ProductID)).Where(q.RuleExecution.RuleID.Eq(ruleID), q.Product.TenantID.Eq(tenant.GetTenantID(ctx)))
	items, total, err := executions.Order(q.RuleExecution.CreatedAt.Desc()).FindByPage((page-1)*size, size)
	if err != nil {
		return nil, 0, err
	}
	result := make([]model.RuleExecution, len(items))
	for i := range items {
		result[i] = *items[i]
	}
	return result, total, nil
}

func (r *RuleRepository) CreateExecution(ctx context.Context, execution *model.RuleExecution) error {
	if _, err := r.Find(ctx, execution.RuleID); err != nil {
		return err
	}
	return useIoTQuery(r.db).RuleExecution.WithContext(ctx).Create(execution)
}

func (r *RuleRepository) UpdateExecutionStats(ctx context.Context, rule *model.Rule, executedAt time.Time) error {
	q := useIoTQuery(r.db)
	_, err := q.Rule.WithContext(ctx).Where(q.Rule.ID.Eq(rule.ID)).UpdateSimple(q.Rule.ExecutionCount.Value(rule.ExecutionCount+1), q.Rule.SuccessCount.Value(rule.SuccessCount+1), q.Rule.LastExecutionStatus.Value("success"), q.Rule.LastExecutedAt.Value(executedAt))
	return err
}
