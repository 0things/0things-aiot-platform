package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"0things-backend/internal/model"
	"0things-backend/internal/repository"
)

var ErrInvalidDeviceEvent = errors.New("invalid device event")

type DeviceEventService struct {
	repo    *repository.DeviceEventRepository
	devices *repository.DeviceRepository
}

func NewDeviceEventService(repo *repository.DeviceEventRepository, devices *repository.DeviceRepository) *DeviceEventService {
	return &DeviceEventService{repo: repo, devices: devices}
}

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
	return s.repo.Create(ctx, &model.DeviceEvent{DeviceID: device.ID, EventType: eventType, EventAt: eventAt, Data: payload})
}

func (s *DeviceEventService) List(ctx context.Context, page, size int, keyword, deviceKey, eventType string, startAt, endAt *time.Time) ([]model.DeviceEvent, int64, error) {
	return s.repo.List(ctx, page, size, keyword, deviceKey, eventType, startAt, endAt)
}
