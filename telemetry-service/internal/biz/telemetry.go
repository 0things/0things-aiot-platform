package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// TelemetryEvent represents a telemetry event message.
type TelemetryEvent struct {
	Type       string         `json:"type"`
	Timestamp  int64          `json:"timestamp"`
	ProductKey string         `json:"product_key,omitempty"`
	DeviceKey  string         `json:"device_key,omitempty"`
	SessionID  string         `json:"session_id,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

// TelemetryEventRepo defines the repository interface for storing telemetry events.
type TelemetryEventRepo interface {
	Save(context.Context, *TelemetryEvent) error
}

// TelemetryUsecase handles telemetry event processing.
type TelemetryUsecase struct {
	repo   TelemetryEventRepo
	logger *log.Helper
}

// NewTelemetryUsecase creates a new telemetry use case.
func NewTelemetryUsecase(repo TelemetryEventRepo, logger log.Logger) *TelemetryUsecase {
	return &TelemetryUsecase{
		repo:   repo,
		logger: log.NewHelper(logger),
	}
}

// ProcessEvent processes a telemetry event (business logic only).
// This method doesn't care where the event came from (Kafka, HTTP, gRPC, etc.)
func (uc *TelemetryUsecase) ProcessEvent(ctx context.Context, event *TelemetryEvent) error {
	// Business logic: validate, enrich, transform
	uc.logger.Infof("Processing telemetry event: type=%s, product_key=%s, device_key=%s",
		event.Type, event.ProductKey, event.DeviceKey)

	// You can add business rules here:
	// - Validate event type
	// - Enrich with additional data
	// - Apply business rules
	// - Call other use cases if needed

	// Save the event to the repository
	if err := uc.repo.Save(ctx, event); err != nil {
		uc.logger.Errorf("Failed to save telemetry event: %v", err)
		return err
	}

	uc.logger.Debugf("Successfully processed telemetry event: %s", event.Type)
	return nil
}
