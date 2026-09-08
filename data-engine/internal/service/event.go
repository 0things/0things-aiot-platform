package service

import (
	"context"

	"0things/pkg/event"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// EventService manages lifecycle (online, offline, heartbeat) and business alarm event processing.
type EventService interface {
	HandleEvent(ctx context.Context, msg event.DeviceMessage) error
}

type eventService struct {
	config *viper.Viper
	logger *zap.Logger
}

// NewEventService creates an instance of EventService.
func NewEventService(config *viper.Viper, logger *zap.Logger) EventService {
	return &eventService{
		config: config,
		logger: logger,
	}
}

// HandleEvent processes incoming device events and audits/dispatches notifications.
func (s *eventService) HandleEvent(ctx context.Context, msg event.DeviceMessage) error {
	s.logger.Info("processing device event",
		zap.String("device_key", msg.DeviceKey),
		zap.String("transport", msg.Transport),
		zap.String("type", msg.MessageType),
		zap.Time("timestamp", msg.Timestamp),
	)

	// Pipeline:
	// 1. If online/offline status event, update device presence and active timestamp.
	// 2. If hardware failure/alarm event, record audit log or trigger notification workflow.
	return nil
}
