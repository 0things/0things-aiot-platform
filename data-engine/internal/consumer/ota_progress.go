package consumer

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/handler"
)

// OTAProgressConsumer handles device OTA progress reporting and state aggregation.
type OTAProgressConsumer struct {
	otaHandler handler.OTAProgressHandlerInterface
}

func NewOTAProgressConsumer(otaHandler handler.OTAProgressHandlerInterface) *OTAProgressConsumer {
	return &OTAProgressConsumer{
		otaHandler: otaHandler,
	}
}

// HandleProgressReport handles device OTA progress report and advances state machine.
func (c *OTAProgressConsumer) HandleProgressReport(ctx context.Context, report *event.OTAUpgradeReport, meta map[string]string) error {
	return c.otaHandler.HandleProgressReport(ctx, report, meta)
}
