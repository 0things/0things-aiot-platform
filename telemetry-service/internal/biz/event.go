package biz

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
)

// Event represents a user event.
type Event struct {
	Type       string         `json:"type"` // e.g., "device.online", "device.offline", "device.alert"
	Timestamp  int64          `json:"timestamp"`
	ProductKey string         `json:"product_key"`
	DeviceKey  string         `json:"device_key"`
	Data       map[string]any `json:"data,omitempty"`
}

// EventRepo defines the repository interface for storing events.
type EventRepo interface {
	Save(context.Context, *Event) error
}

// EventUsecase handles event processing.
type EventUsecase struct {
	repo   EventRepo
	logger *log.Helper
}

// NewEventUsecase creates a new event use case.
func NewEventUsecase(repo EventRepo, logger log.Logger) *EventUsecase {
	return &EventUsecase{
		repo:   repo,
		logger: log.NewHelper(logger),
	}
}

// ProcessEvent processes an event (business logic).
func (uc *EventUsecase) ProcessEvent(ctx context.Context, event *Event) error {
	// Business validation
	if event.ProductKey == "" {
		return fmt.Errorf("product_key is required for events")
	}
	if event.DeviceKey == "" {
		return fmt.Errorf("device_key is required for events")
	}

	uc.logger.Infof("Processing event: type=%s, product_key=%s, device_key=%s",
		event.Type, event.ProductKey, event.DeviceKey)

	// Business logic examples:
	switch event.Type {
	case "device.online":
		uc.logger.Infof("Device online: product_key=%s, device_key=%s", event.ProductKey, event.DeviceKey)
	case "device.offline":
		uc.logger.Infof("Device offline: product_key=%s, device_key=%s", event.ProductKey, event.DeviceKey)
	case "device.alert":
		uc.logger.Infof("Device alert: product_key=%s, device_key=%s", event.ProductKey, event.DeviceKey)
	default:
		uc.logger.Warnf("Unknown event type: %s", event.Type)
	}

	// Save the event
	if err := uc.repo.Save(ctx, event); err != nil {
		uc.logger.Errorf("Failed to save event: %v", err)
		return err
	}

	uc.logger.Debugf("Successfully processed event: %s", event.Type)
	return nil
}
