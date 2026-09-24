package handler

import (
	"context"

	"data-engine/internal/event"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

// EventHandler dispatches device events to the service layer.
type EventHandler struct {
	eventService service.EventService
	logger       *zap.Logger
}

// NewEventHandler creates a new EventHandler instance with injected EventService.
func NewEventHandler(eventService service.EventService, logger *zap.Logger) *EventHandler {
	return &EventHandler{eventService: eventService, logger: logger}
}

// HandleEvent forwards the received device event message to the EventService.
func (h *EventHandler) HandleEvent(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	if msg == nil {
		return nil
	}
	h.logger.Debug("handling device event report",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport.String()),
		zap.String("type", msg.MessageType.String()),
	)
	return h.eventService.HandleEvent(ctx, *msg)
}
