package service

import (
	"context"
	"time"

	"0things-backend/internal/model"
	"0things-backend/internal/repository"
)

type RuleService struct {
	repo *repository.RuleRepository
}

func NewRuleService(repo *repository.RuleRepository) *RuleService {
	return &RuleService{repo: repo}
}

func (s *RuleService) List(ctx context.Context, page, size int, ruleType, status, search string) ([]model.Rule, int64, error) {
	return s.repo.List(ctx, page, size, ruleType, status, search)
}

func (s *RuleService) Get(ctx context.Context, id int64) (*model.Rule, error) {
	return s.repo.Find(ctx, id)
}
func (s *RuleService) Create(ctx context.Context, rule *model.Rule) error {
	return s.repo.Create(ctx, rule)
}
func (s *RuleService) Update(ctx context.Context, rule *model.Rule) error {
	return s.repo.Save(ctx, rule)
}
func (s *RuleService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
func (s *RuleService) SetStatus(ctx context.Context, id int64, status string) (*model.Rule, error) {
	rule, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err = s.repo.UpdateStatus(ctx, rule, status); err != nil {
		return nil, err
	}
	rule.Status = status
	return rule, nil
}
func (s *RuleService) ListExecutions(ctx context.Context, id int64, page, size int) ([]model.RuleExecution, int64, error) {
	return s.repo.ListExecutions(ctx, id, page, size)
}

func (s *RuleService) Evaluate(ctx context.Context, id int64) (*model.RuleExecution, error) {
	rule, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	execution := model.RuleExecution{RuleID: id, RuleName: rule.Name, Status: "success", ConditionResult: true, TriggeredAt: time.Now()}
	if err = s.repo.CreateExecution(ctx, &execution); err != nil {
		return nil, err
	}
	if err = s.repo.UpdateExecutionStats(ctx, rule, execution.TriggeredAt); err != nil {
		return nil, err
	}
	return &execution, nil
}
