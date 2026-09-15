package handler

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

// OTAHandler dispatches OTA progress and inform reports to the service layer.
type OTAHandler struct {
	otaService service.OTAService
	logger     *zap.Logger
}

// NewOTAHandler creates a new OTAHandler instance with injected OTAService.
func NewOTAHandler(otaService service.OTAService, logger *zap.Logger) *OTAHandler {
	return &OTAHandler{otaService: otaService, logger: logger}
}

// HandleProgressReport processes incoming OTA progress reports and dispatches them to OTAService.
func (h *OTAHandler) HandleProgressReport(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	if msg == nil {
		return nil
	}
	h.logger.Info("handling OTA progress report event",
		zap.String("device_key", msg.DeviceKey),
		zap.Int("payload_bytes", len(msg.Payload)),
	)
	return h.otaService.HandleProgressReport(ctx, *msg)
}

// HandleDeviceInfo processes incoming OTA inform reports and dispatches them to OTAService.
func (h *OTAHandler) HandleDeviceInfo(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	if msg == nil {
		return nil
	}
	h.logger.Info("handling OTA device info report event",
		zap.String("device_key", msg.DeviceKey),
		zap.Int("payload_bytes", len(msg.Payload)),
	)
	return h.otaService.HandleDeviceInfo(ctx, *msg)
}
