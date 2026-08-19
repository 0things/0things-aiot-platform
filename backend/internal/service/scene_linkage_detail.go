package service

import (
	"context"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
)

type SceneLinkageDetailServiceInterface interface {
	GetBySceneID(ctx context.Context, sceneID int64) (*model.SceneLinkageDetail, error)
	Create(ctx context.Context, detail *model.SceneLinkageDetail) error
	Update(ctx context.Context, detail *model.SceneLinkageDetail) error
}

type SceneLinkageDetailService struct {
	repo *repository.SceneLinkageDetailRepository
}

func NewSceneLinkageDetailService(repo *repository.SceneLinkageDetailRepository) *SceneLinkageDetailService {
	return &SceneLinkageDetailService{repo: repo}
}

func (s *SceneLinkageDetailService) GetBySceneID(ctx context.Context, sceneID int64) (*model.SceneLinkageDetail, error) {
	return s.repo.FindBySceneID(ctx, sceneID)
}

func (s *SceneLinkageDetailService) Create(ctx context.Context, detail *model.SceneLinkageDetail) error {
	return s.repo.Create(ctx, detail)
}

func (s *SceneLinkageDetailService) Update(ctx context.Context, detail *model.SceneLinkageDetail) error {
	return s.repo.Save(ctx, detail)
}
