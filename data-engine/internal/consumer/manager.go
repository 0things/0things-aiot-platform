package consumer

import (
	"context"

	"0things/pkg/event"

	"go.uber.org/zap"
)

// Manager coordinates event subscriptions across all domain consumer handlers.
type Manager struct {
	consumer            event.Consumer
	telemetryConsumer   *TelemetryConsumer
	eventConsumer       *EventConsumer
	otaProgressConsumer *OTAProgressConsumer
	logger              *zap.Logger
}

func NewManager(
	consumer event.Consumer,
	telemetryConsumer *TelemetryConsumer,
	eventConsumer *EventConsumer,
	otaProgressConsumer *OTAProgressConsumer,
	logger *zap.Logger,
) *Manager {
	return &Manager{
		consumer:            consumer,
		telemetryConsumer:   telemetryConsumer,
		eventConsumer:       eventConsumer,
		otaProgressConsumer: otaProgressConsumer,
		logger:              logger,
	}
}

// Start registers all topic subscriptions with the event bus.
func (m *Manager) Start(ctx context.Context) error {
	// 1. Subscribe to device telemetry reports
	if err := event.Subscribe(ctx, m.consumer, event.TopicDeviceTelemetryReport, m.telemetryConsumer.HandleTelemetry); err != nil {
		return err
	}

	// 2. Subscribe to device attribute reports
	if err := event.Subscribe(ctx, m.consumer, event.TopicDeviceAttributeReport, m.telemetryConsumer.HandleAttribute); err != nil {
		return err
	}

	// 3. Subscribe to device business events and alarms
	if err := event.Subscribe(ctx, m.consumer, event.TopicDeviceEventReport, m.eventConsumer.HandleEvent); err != nil {
		return err
	}

	// 4. Subscribe to device OTA progress reports
	if err := event.Subscribe(ctx, m.consumer, event.TopicOTAProgressReport, m.otaProgressConsumer.HandleProgressReport); err != nil {
		return err
	}

	m.logger.Info("all event subscriptions registered successfully")
	return nil
}
