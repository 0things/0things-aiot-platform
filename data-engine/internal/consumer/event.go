package consumer

import (
	"context"

	"data-engine/internal/event"
	"data-engine/internal/handler"
)

// EventConsumer handles device lifecycle and business alarm events.
type EventConsumer struct {
	eventHandler *handler.EventHandler
}

// NewEventConsumer creates a new EventConsumer instance with injected EventHandler.
func NewEventConsumer(eventHandler *handler.EventHandler) *EventConsumer {
	return &EventConsumer{
		eventHandler: eventHandler,
	}
}

// HandleEvent processes device event reports and forwards them to the handler.
func (c *EventConsumer) HandleEvent(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	return c.eventHandler.HandleEvent(ctx, msg, meta)
}
