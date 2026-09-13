package handler

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

// OTAProgressHandlerInterface handles OTA progress reports.
type OTAProgressHandlerInterface interface {
	HandleProgressReport(context.Context, *event.OTAUpgradeReport, map[string]string) error
}

// OTAProgressHandler dispatches OTA reports to the service layer.
type OTAProgressHandler struct {
	otaService service.OTAService
	logger     *zap.Logger
}

func NewOTAProgressHandler(otaService service.OTAService, logger *zap.Logger) *OTAProgressHandler {
	return &OTAProgressHandler{otaService: otaService, logger: logger}
}

func (h *OTAProgressHandler) HandleProgressReport(ctx context.Context, report *event.OTAUpgradeReport, meta map[string]string) error {
	if report == nil {
		return nil
	}
	h.logger.Info("handling OTA progress report event",
		zap.String("device_key", report.DeviceKey),
		zap.String("batch_id", report.BatchID),
		zap.String("status", report.Status),
	)

	rep := *report
	if rep.EventType == "" {
		rep.EventType = "progress"
	}
	return h.otaService.HandleOTAReport(ctx, rep)
}
