package consumer

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

// EventConsumer handles device lifecycle and business alarm events.
type EventConsumer struct {
	eventService service.EventService
	logger       *zap.Logger
}

func NewEventConsumer(eventService service.EventService, logger *zap.Logger) *EventConsumer {
	return &EventConsumer{
		eventService: eventService,
		logger:       logger,
	}
}

// HandleEvent processes device event reports.
func (c *EventConsumer) HandleEvent(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	if msg == nil {
		return nil
	}
	c.logger.Debug("handling device event report",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport),
		zap.String("type", msg.MessageType),
	)

	return c.eventService.HandleEvent(ctx, *msg)
}
