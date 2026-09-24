package consumer

import (
	"context"

	"data-engine/internal/event"
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

// HandlePropertyPost processes device property/telemetry uplink events and forwards them to the handler.
func (c *TelemetryConsumer) HandlePropertyPost(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	return c.telemetryHandler.HandlePropertyPost(ctx, msg, meta)
}
