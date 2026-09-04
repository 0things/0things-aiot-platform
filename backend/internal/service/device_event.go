package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	eventV1 "aiot-backend/api/v1"
	"aiot-backend/internal/dto"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"github.com/google/uuid"
)

var ErrInvalidDeviceEvent = errors.New("invalid device event")

type DeviceEventServiceInterface interface {
	Record(ctx context.Context, productKey, deviceKey, eventType string, timestamp int64, data map[string]any) error
	List(ctx context.Context, req *eventV1.ListDeviceEventsRequest) ([]dto.DeviceEventListItem, int64, error)
}

type DeviceEventService struct {
	repo    *repository.DeviceEventRepository
	devices *repository.DeviceRepository
}

func NewDeviceEventService(repo *repository.DeviceEventRepository, devices *repository.DeviceRepository) *DeviceEventService {
	return &DeviceEventService{repo: repo, devices: devices}
}

// Record persists a device event after validating device ownership and payload formatting.
func (s *DeviceEventService) Record(ctx context.Context, productKey, deviceKey, eventType string, timestamp int64, data map[string]any) error {
	if productKey == "" || deviceKey == "" || eventType == "" {
		return fmt.Errorf("%w: product_key, device_key and type are required", ErrInvalidDeviceEvent)
	}
	device, err := s.devices.FindByKeyForEvent(ctx, deviceKey)
	if err != nil {
		return err
	}
	if device.Product.ProductKey != productKey {
		return fmt.Errorf("%w: product_key does not match device", ErrInvalidDeviceEvent)
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	eventAt := time.Now()
	if timestamp > 0 {
		eventAt = time.UnixMilli(timestamp)
	}
	// Each event needs an independent reference, even when the same device emits
	// the same event type repeatedly.
	return s.repo.Create(ctx, &model.DeviceEvent{
		UUID:            uuid.NewString(),
		DeviceID:        device.ID,
		EventIdentifier: eventType,
		EventType:       eventType,
		EventAt:         eventAt,
		Data:            string(payload),
	})
}

// List resolves API request parameters into internal query DTO and queries device events.
func (s *DeviceEventService) List(ctx context.Context, req *eventV1.ListDeviceEventsRequest) ([]dto.DeviceEventListItem, int64, error) {
	if req == nil {
		return nil, 0, errors.New("request cannot be nil")
	}

	query := dto.ListDeviceEventsQuery{
		Page:      req.Page,
		PageSize:  req.PageSize,
		Keyword:   req.Keyword,
		DeviceKey: req.DeviceKey,
		EventType: req.EventType,
		StartAt:   req.StartAt,
		EndAt:     req.EndAt,
	}

	return s.repo.List(ctx, query)
}
