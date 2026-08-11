package service

import (
	"context"
	"time"

	"0things-backend/internal/model"
	"0things-backend/internal/repository"
)

type AlertService struct{ repo *repository.AlertRepository }

func NewAlertService(repo *repository.AlertRepository) *AlertService {
	return &AlertService{repo: repo}
}
func (s *AlertService) List(ctx context.Context, page, size int, status, severity, deviceKey string) ([]model.Alert, int64, error) {
	return s.repo.List(ctx, page, size, status, severity, deviceKey)
}

func (s *AlertService) Get(ctx context.Context, id int64) (*model.Alert, error) {
	return s.repo.Find(ctx, id)
}
func (s *AlertService) SetStatus(ctx context.Context, id int64, status string) (*model.Alert, error) {
	alert, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	if err = s.repo.UpdateStatus(ctx, alert, status, now); err != nil {
		return nil, err
	}
	alert.Status = status
	if status == "acknowledged" {
		alert.AckAt = &now
	} else {
		alert.ResolvedAt = &now
	}
	return alert, nil
}
