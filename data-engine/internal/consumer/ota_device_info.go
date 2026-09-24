package consumer

import (
	"context"

	"data-engine/internal/event"
	"data-engine/internal/handler"
)

// OTADeviceInfoConsumer handles device OTA inform version reports and validates upgrade success.
type OTADeviceInfoConsumer struct {
	otaHandler *handler.OTAHandler
}

// NewOTADeviceInfoConsumer creates a new OTADeviceInfoConsumer instance with injected OTAHandler.
func NewOTADeviceInfoConsumer(otaHandler *handler.OTAHandler) *OTADeviceInfoConsumer {
	return &OTADeviceInfoConsumer{
		otaHandler: otaHandler,
	}
}

// HandleDeviceInfo handles device OTA inform version report and validates success.
func (c *OTADeviceInfoConsumer) HandleDeviceInfo(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	return c.otaHandler.HandleDeviceInfo(ctx, msg, meta)
}
