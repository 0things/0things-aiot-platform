package consumer

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

// OTAProgressConsumer handles device OTA progress reporting and state aggregation.
type OTAProgressConsumer struct {
	processor *service.OTAProcessor
	logger    *zap.Logger
}

func NewOTAProgressConsumer(processor *service.OTAProcessor, logger *zap.Logger) *OTAProgressConsumer {
	return &OTAProgressConsumer{
		processor: processor,
		logger:    logger,
	}
}

// HandleProgressReport handles device OTA progress report and advances state machine.
func (c *OTAProgressConsumer) HandleProgressReport(ctx context.Context, report *event.OTAUpgradeReport, meta map[string]string) error {
	c.logger.Info("handling OTA progress report event",
		zap.String("device_key", report.DeviceKey),
		zap.String("batch_id", report.BatchID),
		zap.String("status", report.Status),
	)

	rep := *report
	if rep.EventType == "" {
		rep.EventType = "progress"
	}

	return c.processor.HandleOTAReport(ctx, rep)
}
