package service

import (
	"context"

	"0things-backend/internal/model"
	"0things-backend/internal/repository"
)

type SceneLinkageServiceInterface interface {
	List(ctx context.Context, page, size int, search string, enable int) ([]model.SceneLinkage, int64, error)
	Get(ctx context.Context, id int64) (*model.SceneLinkage, error)
	Create(ctx context.Context, scene *model.SceneLinkage) error
	Update(ctx context.Context, scene *model.SceneLinkage) error
	Delete(ctx context.Context, id int64) error
}

type SceneLinkageService struct {
	repo *repository.SceneLinkageRepository
}

func NewSceneLinkageService(repo *repository.SceneLinkageRepository) *SceneLinkageService {
	return &SceneLinkageService{repo: repo}
}

func (s *SceneLinkageService) List(ctx context.Context, page, size int, search string, enable int) ([]model.SceneLinkage, int64, error) {
	return s.repo.List(ctx, page, size, search, enable)
}

func (s *SceneLinkageService) Get(ctx context.Context, id int64) (*model.SceneLinkage, error) {
	return s.repo.Find(ctx, id)
}

func (s *SceneLinkageService) Create(ctx context.Context, scene *model.SceneLinkage) error {
	return s.repo.Create(ctx, scene)
}

func (s *SceneLinkageService) Update(ctx context.Context, scene *model.SceneLinkage) error {
	return s.repo.Save(ctx, scene)
}

func (s *SceneLinkageService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
