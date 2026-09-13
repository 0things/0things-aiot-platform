package consumer

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/handler"
)

// EventConsumer handles device lifecycle and business alarm events.
type EventConsumer struct {
	eventHandler handler.EventHandlerInterface
}

func NewEventConsumer(eventHandler handler.EventHandlerInterface) *EventConsumer {
	return &EventConsumer{
		eventHandler: eventHandler,
	}
}

// HandleEvent processes device event reports.
func (c *EventConsumer) HandleEvent(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	return c.eventHandler.HandleEvent(ctx, msg, meta)
}
