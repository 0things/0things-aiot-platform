package handler

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

// TelemetryHandler dispatches telemetry-related events to the service layer.
type TelemetryHandler struct {
	telemetryService service.TelemetryService
	logger           *zap.Logger
}

// NewTelemetryHandler creates a new TelemetryHandler instance with injected TelemetryService.
func NewTelemetryHandler(telemetryService service.TelemetryService, logger *zap.Logger) *TelemetryHandler {
	return &TelemetryHandler{telemetryService: telemetryService, logger: logger}
}

// HandlePropertyPost forwards device property/telemetry metrics to TelemetryService for processing.
func (h *TelemetryHandler) HandlePropertyPost(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	if msg == nil {
		return nil
	}
	h.logger.Debug("handling device property post report from event bus",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport.String()),
	)
	return h.telemetryService.ProcessPropertyPost(ctx, *msg)
}
