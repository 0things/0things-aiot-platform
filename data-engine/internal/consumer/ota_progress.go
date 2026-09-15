package consumer

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/handler"
)

// OTAProgressConsumer handles device OTA progress reporting and state aggregation.
type OTAProgressConsumer struct {
	otaHandler *handler.OTAHandler
}

// NewOTAProgressConsumer creates a new OTAProgressConsumer instance with injected OTAHandler.
func NewOTAProgressConsumer(otaHandler *handler.OTAHandler) *OTAProgressConsumer {
	return &OTAProgressConsumer{
		otaHandler: otaHandler,
	}
}

// HandleProgressReport handles device OTA progress report and advances state machine.
func (c *OTAProgressConsumer) HandleProgressReport(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	return c.otaHandler.HandleProgressReport(ctx, msg, meta)
}
