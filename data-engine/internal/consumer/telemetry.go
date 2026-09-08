package consumer

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

// TelemetryConsumer handles device telemetry and attribute reports.
type TelemetryConsumer struct {
	telemetryService service.TelemetryService
	logger           *zap.Logger
}

func NewTelemetryConsumer(telemetryService service.TelemetryService, logger *zap.Logger) *TelemetryConsumer {
	return &TelemetryConsumer{
		telemetryService: telemetryService,
		logger:           logger,
	}
}

// HandleTelemetry processes device telemetry uplink events.
func (c *TelemetryConsumer) HandleTelemetry(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	if msg == nil {
		return nil
	}
	c.logger.Debug("handling device telemetry report",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport),
	)

	return c.telemetryService.ProcessMessage(ctx, *msg)
}

// HandleAttribute processes device attribute update events.
func (c *TelemetryConsumer) HandleAttribute(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	if msg == nil {
		return nil
	}
	c.logger.Debug("handling device attribute report",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport),
	)

	return c.telemetryService.ProcessMessage(ctx, *msg)
}
