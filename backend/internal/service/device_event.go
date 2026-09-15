package service

import (
	"context"

	"aiot-backend/internal/dto"
	"aiot-backend/internal/repository"
)

// DeviceEventServiceInterface defines service operations for managing device events.
type DeviceEventServiceInterface interface {
	// List queries paginated device events matching the query criteria.
	List(ctx context.Context, query dto.ListDeviceEventsQuery) ([]dto.DeviceEventListItem, int64, error)
}

type DeviceEventService struct {
	repo *repository.DeviceEventRepository
}

// NewDeviceEventService creates a new DeviceEventService instance.
func NewDeviceEventService(repo *repository.DeviceEventRepository) *DeviceEventService {
	return &DeviceEventService{repo: repo}
}

// List queries device events by internal query DTO.
func (s *DeviceEventService) List(ctx context.Context, query dto.ListDeviceEventsQuery) ([]dto.DeviceEventListItem, int64, error) {
	return s.repo.List(ctx, query)
}
