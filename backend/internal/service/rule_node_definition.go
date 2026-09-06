package service

import (
	"context"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
)

type RuleNodeDefinitionServiceInterface interface {
	ListEnabled(ctx context.Context) ([]model.RuleNodeDefinition, error)
}

type RuleNodeDefinitionService struct {
	repo *repository.RuleNodeDefinitionRepository
}

func NewRuleNodeDefinitionService(repo *repository.RuleNodeDefinitionRepository) *RuleNodeDefinitionService {
	return &RuleNodeDefinitionService{repo: repo}
}

func (s *RuleNodeDefinitionService) ListEnabled(ctx context.Context) ([]model.RuleNodeDefinition, error) {
	return s.repo.ListEnabled(ctx)
}
