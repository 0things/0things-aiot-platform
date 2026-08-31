package service

import (
	"context"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/pkg/log"
)

type TelemetryServiceInterface interface {
	QueryHistory(ctx context.Context, req model.TelemetryQueryReq) ([]model.TelemetryPoint, error)
}

type TelemetryService struct {
	repo   *repository.TelemetryRepository
	logger *log.Logger
}

func NewTelemetryService(repo *repository.TelemetryRepository, logger *log.Logger) *TelemetryService {
	return &TelemetryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *TelemetryService) QueryHistory(ctx context.Context, req model.TelemetryQueryReq) ([]model.TelemetryPoint, error) {
	return s.repo.QueryHistory(ctx, req)
}
