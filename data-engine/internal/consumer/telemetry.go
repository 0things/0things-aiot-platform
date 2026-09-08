package consumer

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/engine"

	"go.uber.org/zap"
)

// TelemetryConsumer handles device telemetry and attribute reports.
type TelemetryConsumer struct {
	processor *engine.Processor
	logger    *zap.Logger
}

func NewTelemetryConsumer(processor *engine.Processor, logger *zap.Logger) *TelemetryConsumer {
	return &TelemetryConsumer{
		processor: processor,
		logger:    logger,
	}
}

// HandleTelemetry processes device telemetry uplink events.
func (c *TelemetryConsumer) HandleTelemetry(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	c.logger.Debug("handling device telemetry report",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport),
	)

	if msg == nil {
		return nil
	}
	return c.processor.ProcessMessage(ctx, *msg)
}

// HandleAttribute processes device attribute update events.
func (c *TelemetryConsumer) HandleAttribute(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	c.logger.Debug("handling device attribute report",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport),
	)

	if msg == nil {
		return nil
	}
	return c.processor.ProcessMessage(ctx, *msg)
}
