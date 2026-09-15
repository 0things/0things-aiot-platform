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

// HandleTelemetry forwards device telemetry metrics to TelemetryService for processing.
func (h *TelemetryHandler) HandleTelemetry(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	if msg == nil {
		return nil
	}
	h.logger.Debug("handling device telemetry report",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport.String()),
	)
	return h.telemetryService.ProcessMessage(ctx, *msg)
}

// HandleAttribute forwards device attribute updates to TelemetryService for processing.
func (h *TelemetryHandler) HandleAttribute(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	if msg == nil {
		return nil
	}
	h.logger.Debug("handling device attribute report",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport.String()),
	)
	return h.telemetryService.ProcessMessage(ctx, *msg)
}
