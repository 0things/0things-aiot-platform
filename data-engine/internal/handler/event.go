package handler

import (
	"context"

	"0things/pkg/event"
	"data-engine/internal/service"

	"go.uber.org/zap"
)

// EventHandlerInterface handles device lifecycle and business alarm events.
type EventHandlerInterface interface {
	HandleEvent(context.Context, *event.DeviceMessage, map[string]string) error
}

// EventHandler dispatches device events to the service layer.
type EventHandler struct {
	eventService service.EventService
	logger       *zap.Logger
}

func NewEventHandler(eventService service.EventService, logger *zap.Logger) *EventHandler {
	return &EventHandler{eventService: eventService, logger: logger}
}

func (h *EventHandler) HandleEvent(ctx context.Context, msg *event.DeviceMessage, meta map[string]string) error {
	if msg == nil {
		return nil
	}
	h.logger.Debug("handling device event report",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport),
		zap.String("type", msg.MessageType),
	)
	return h.eventService.HandleEvent(ctx, *msg)
}
