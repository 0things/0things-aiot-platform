package handler

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

// TelemetryHandlerInterface handles telemetry and attribute reports.
type TelemetryHandlerInterface interface {
	HandleTelemetry(context.Context, *event.DeviceMessage, map[string]string) error
	HandleAttribute(context.Context, *event.DeviceMessage, map[string]string) error
}

// TelemetryHandler dispatches telemetry-related events to the service layer.
type TelemetryHandler struct {
	telemetryService service.TelemetryService
	logger           *zap.Logger
}

func NewTelemetryHandler(telemetryService service.TelemetryService, logger *zap.Logger) *TelemetryHandler {
	return &TelemetryHandler{telemetryService: telemetryService, logger: logger}
}

func (h *TelemetryHandler) HandleTelemetry(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	if msg == nil {
		return nil
	}
	h.logger.Debug("handling device telemetry report",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport),
	)
	return h.telemetryService.ProcessMessage(ctx, *msg)
}

func (h *TelemetryHandler) HandleAttribute(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	if msg == nil {
		return nil
	}
	h.logger.Debug("handling device attribute report",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport),
	)
	return h.telemetryService.ProcessMessage(ctx, *msg)
}
