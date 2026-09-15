package consumer

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/handler"
)

// TelemetryConsumer handles device telemetry and attribute reports.
type TelemetryConsumer struct {
	telemetryHandler *handler.TelemetryHandler
}

// NewTelemetryConsumer creates a new TelemetryConsumer instance with injected TelemetryHandler.
func NewTelemetryConsumer(telemetryHandler *handler.TelemetryHandler) *TelemetryConsumer {
	return &TelemetryConsumer{
		telemetryHandler: telemetryHandler,
	}
}

// HandleTelemetry processes device telemetry uplink events and forwards them to the handler.
func (c *TelemetryConsumer) HandleTelemetry(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	return c.telemetryHandler.HandleTelemetry(ctx, msg, meta)
}

// HandleAttribute processes device attribute update events and forwards them to the handler.
func (c *TelemetryConsumer) HandleAttribute(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	return c.telemetryHandler.HandleAttribute(ctx, msg, meta)
}
